package usecase

import (
	"backend/internal/audit"
	"backend/internal/cache"
	"backend/internal/domain"
	pb "backend/internal/proto/hiring/v1"
	"backend/internal/repo"
	"backend/pkg/config"
	"backend/pkg/hash"
	"backend/pkg/mq"
	"backend/pkg/rbac"
	"backend/pkg/token"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInviteNotFound = errors.New("invite not found")
	ErrInviteExpired  = errors.New("invite expired")
)

const inviteTTL = 48 * time.Hour

type InviteUseCase interface {
	InviteUser(ctx context.Context, tokenStr string, req domain.CreateInviteParams) error
	ValidateInvite(ctx context.Context, token string) (*domain.InviteRegisterDTO, error)
	AcceptInvite(ctx context.Context, req domain.CreateUserParams) (*string, time.Duration, error)
	CompleteInvite(ctx context.Context, inviteID string, userID string, jobIDs []string) error
	ProcessStuckInvites(ctx context.Context) error
	GetTeamInvites(ctx context.Context, teamID string) ([]domain.Invite, error)
}

var _ InviteUseCase = (*inviteUseCase)(nil)

type inviteUseCase struct {
	cfg          *config.Config
	repo         repo.InviteRepository
	userRepo     repo.UserRepository
	cacheManager *cache.Manager
	token        *token.JWTtoken
	hash         hash.Hash
	enforcer     *rbac.CasbinClient
	hiringClient *pb.HiringServiceClient
	auditor      *audit.Logger
	publisher    *mq.MQPublisher
}

func NewInviteUseCase(
	cfg *config.Config,
	repo repo.InviteRepository,
	userRepo repo.UserRepository,
	cacheManager *cache.Manager,
	token *token.JWTtoken,
	hash hash.Hash,
	enforcer *rbac.CasbinClient,
	hiringClient *pb.HiringServiceClient,
	auditor *audit.Logger,
	publisher *mq.MQPublisher,
) InviteUseCase {
	return &inviteUseCase{
		cfg:          cfg,
		repo:         repo,
		userRepo:     userRepo,
		cacheManager: cacheManager,
		token:        token,
		hash:         hash,
		enforcer:     enforcer,
		hiringClient: hiringClient,
		auditor:      auditor,
		publisher:    publisher,
	}
}

func (i *inviteUseCase) InviteUser(ctx context.Context, tokenStr string, req domain.CreateInviteParams) error {
	token, err := i.token.VerifyToken([]byte(tokenStr))
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	subject, err := cache.Get(ctx, i.cacheManager, cache.SessionKey, token.Subject())
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}

	inviteToken := uuid.New().String()

	invite := &domain.Invite{
		TeamID:    subject.TeamID,
		Email:     req.Email,
		Role:      req.Role,
		Token:     inviteToken,
		ExpiresAt: time.Now().Add(i.cfg.Invite.TTL),
	}

	var jobIDs []string
	if req.JobIDs != nil {
		jobIDs = *req.JobIDs
	}

	if err := i.repo.CreateInvite(ctx, invite, jobIDs); err != nil {
		return fmt.Errorf("create invite: %w", err)
	}

	_ = i.auditor.Log(ctx, audit.Entry{
		TeamID:    subject.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &subject.UserID,
		Action:    audit.AuthInviteSent,
		TargetID:  &invite.ID,
		Payload:   map[string]string{"email": req.Email, "role": req.Role},
	})

	payload, err := json.Marshal(domain.InviteCreatedEvent{
		InviteID: invite.ID,
		Email:    invite.Email,
		Token:    inviteToken,
		TeamID:   invite.TeamID,
		Role:     invite.Role,
		Locale:   req.Locale,
	})
	if err != nil {
		return fmt.Errorf("marshal invite created event: %w", err)
	}

	if err := i.publisher.Publish(ctx, payload,
		mq.WithExchange("hiring.events"),
		mq.WithRoutingKey("invite.created"),
	); err != nil {
		return fmt.Errorf("publish invite created event: %w", err)
	}

	return nil
}

func (i *inviteUseCase) ValidateInvite(ctx context.Context, token string) (*domain.InviteRegisterDTO, error) {
	invite, err := i.repo.GetInviteByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repo.ErrInviteNotFound) {
			return nil, ErrInviteNotFound
		}

		return nil, fmt.Errorf("get invite: %w", err)
	}

	if time.Now().After(invite.ExpiresAt) {
		return nil, ErrInviteExpired
	}

	return &domain.InviteRegisterDTO{Email: invite.Email}, nil
}

func (i *inviteUseCase) AcceptInvite(ctx context.Context, req domain.CreateUserParams) (*string, time.Duration, error) {
	invite, jobIDs, err := i.repo.LockInviteForProcessing(ctx, req.Token)
	if err != nil {
		if errors.Is(err, repo.ErrInviteNotFound) {
			return nil, 0, ErrInviteNotFound
		}
		return nil, 0, fmt.Errorf("lock invite for processing: %w", err)
	}

	if time.Now().After(invite.ExpiresAt) {
		if err := i.repo.UpdateInviteStatus(ctx, invite.ID, "failed"); err != nil {
			return nil, 0, fmt.Errorf("update invite status: %w", err)
		}

		return nil, 0, ErrInviteExpired
	}

	hashedPassword, err := i.hash.Hash(req.Password)
	if err != nil {
		if err := i.repo.RollbackInviteToPending(ctx, invite.ID); err != nil {
			return nil, 0, fmt.Errorf("rollback invite to pending: %w", err)
		}

		return nil, 0, fmt.Errorf("hash password: %w", err)
	}

	user, err := i.repo.AcceptInviteLocalTx(ctx, &domain.CreateUserRepoParams{
		Email:     invite.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  hashedPassword,
		TeamID:    invite.TeamID,
		Role:      invite.Role,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if err := i.repo.UpdateInviteStatus(ctx, invite.ID, "failed"); err != nil {
				return nil, 0, fmt.Errorf("update invite status: %w", err)
			}

			return nil, 0, ErrUserAlreadyExists
		}

		if err := i.repo.RollbackInviteToPending(ctx, invite.ID); err != nil {
			return nil, 0, fmt.Errorf("rollback invite to pending: %w", err)
		}

		return nil, 0, fmt.Errorf("accept invite local tx: %w", err)
	}

	if _, err := i.enforcer.AddRoleForUserInDomain(user.ID, invite.Role, invite.TeamID); err != nil {
		return nil, 0, fmt.Errorf("add role for user in domain: %w", err)
	}

	sessionID := uuid.New().String()
	if err := cache.SetWithTTL(ctx, i.cacheManager, cache.SessionKey, sessionID, domain.Session{
		UserID: user.ID,
		TeamID: user.TeamID,
		Role:   user.Role,
	}, i.token.ExpireAt); err != nil {
		return nil, 0, fmt.Errorf("set session: %w", err)
	}

	signed, err := i.token.GenerateToken(sessionID)
	if err != nil {
		return nil, 0, fmt.Errorf("generate token: %w", err)
	}
	tokenString := string(signed)

	if err := i.CompleteInvite(ctx, invite.ID, user.ID, jobIDs); err != nil {
		return &tokenString, i.token.ExpireAt, nil
	}

	_ = i.auditor.Log(ctx, audit.Entry{
		TeamID:    invite.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &user.ID,
		Action:    audit.AuthInviteAccepted,
		TargetID:  &invite.ID,
	})

	return &tokenString, i.token.ExpireAt, nil
}

func (i *inviteUseCase) CompleteInvite(ctx context.Context, inviteID string, userID string, jobIDs []string) error {
	if len(jobIDs) > 0 {
		grpcReq := &pb.TransferJobAccessRequest{
			UserId: userID,
			JobIds: jobIDs,
		}

		resp, err := (*i.hiringClient).TransferJobAccess(ctx, grpcReq)
		if err != nil {
			return fmt.Errorf("transfer job access: %w", err)
		}

		if !resp.Success {
			return fmt.Errorf("transfer job access failed: unsuccessful response")
		}
	}

	if err := i.repo.UpdateInviteStatus(ctx, inviteID, "completed"); err != nil {
		return fmt.Errorf("update invite status: %w", err)
	}

	return nil
}

func (i *inviteUseCase) ProcessStuckInvites(ctx context.Context) error {
	invites, err := i.repo.GetStuckInvites(ctx, i.cfg.InviteRecovery.StuckThreshold)
	if err != nil {
		return fmt.Errorf("get stuck invites: %w", err)
	}

	for _, invite := range invites {
		user, err := i.userRepo.Login(ctx, invite.Email)
		if err != nil {
			continue
		}

		jobIDs, err := i.repo.GetInviteJobIDs(ctx, invite.ID)
		if err != nil {
			continue
		}

		_ = i.CompleteInvite(ctx, invite.ID, user.ID, jobIDs)
	}

	return nil
}

func (i *inviteUseCase) GetTeamInvites(ctx context.Context, teamID string) ([]domain.Invite, error) {
	invites, err := i.repo.GetInvitesByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("repo get team invites: %w", err)
	}
	return invites, nil
}
