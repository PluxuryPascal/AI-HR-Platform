package usecase

import (
	"backend/internal/cache"
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/pkg/config"
	"backend/pkg/hash"
	"backend/pkg/token"
	"context"
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
}

var _ InviteUseCase = (*inviteUseCase)(nil)

type inviteUseCase struct {
	cfg          *config.Config
	repo         repo.InviteRepository
	cacheManager *cache.Manager
	token        *token.JWTtoken
	hash         hash.Hash
}

func NewInviteUseCase(
	cfg *config.Config,
	repo repo.InviteRepository,
	cacheManager *cache.Manager,
	token *token.JWTtoken,
	hash hash.Hash,
) InviteUseCase {
	return &inviteUseCase{
		cfg:          cfg,
		repo:         repo,
		cacheManager: cacheManager,
		token:        token,
		hash:         hash,
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

	// TODO: hand the token to an email-sending service here.
	// e.g. emailSvc.SendInvite(ctx, invite.Email, inviteToken)

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
	invite, err := i.repo.GetInviteByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, repo.ErrInviteNotFound) {
			return nil, 0, ErrInviteNotFound
		}

		return nil, 0, fmt.Errorf("get invite by token: %w", err)
	}

	if time.Now().After(invite.ExpiresAt) {
		return nil, 0, ErrInviteExpired
	}

	hashedPassword, err := i.hash.Hash(req.Password)
	if err != nil {
		return nil, 0, fmt.Errorf("hash password: %w", err)
	}

	user, err := i.repo.AcceptInviteAndCreateUser(ctx, invite.ID, &domain.CreateUserRepoParams{
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
			return nil, 0, ErrUserAlreadyExists
		}

		return nil, 0, fmt.Errorf("accept invite: %w", err)
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

	return &tokenString, i.token.ExpireAt, nil
}
