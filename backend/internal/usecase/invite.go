package usecase

import (
	"backend/internal/cache"
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/pkg/hash"
	"backend/pkg/token"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type InviteUseCase interface {
	InviteUser(ctx context.Context, req domain.InvitePayload) error
	ValidateInvite(ctx context.Context, req domain.InvitePayload) error
	CreateUser(ctx context.Context, req domain.CreateUserParams) (*string, time.Duration, error)
}

var _ InviteUseCase = (*inviteUseCase)(nil)

type inviteUseCase struct {
	repo         repo.InviteRepository
	cacheManager *cache.Manager
	token        *token.JWTtoken
	hash         hash.Hash
}

func NewInviteUseCase(repo repo.InviteRepository, cacheManager *cache.Manager, token *token.JWTtoken, hash hash.Hash) InviteUseCase {
	return &inviteUseCase{
		repo:         repo,
		cacheManager: cacheManager,
		token:        token,
		hash:         hash,
	}
}

func (i *inviteUseCase) CreateUser(ctx context.Context, req domain.CreateUserParams) (*string, time.Duration, error) {
	invite, err := cache.Get(ctx, i.cacheManager, cache.InviteKey, req.Token)
	if err != nil {
		return nil, 0, fmt.Errorf("get invite from cache: %w", err)
	}

	hashedPassword, err := i.hash.Hash(req.Password)
	if err != nil {
		return nil, 0, fmt.Errorf("hash password: %w", err)
	}

	user, err := i.repo.CreateUser(ctx, &domain.CreateUserRepoParams{
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

		return nil, 0, fmt.Errorf("create user: %w", err)
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

func (i *inviteUseCase) InviteUser(ctx context.Context, req domain.InvitePayload) error {
	panic("unimplemented")
}

func (i *inviteUseCase) ValidateInvite(ctx context.Context, req domain.InvitePayload) error {
	panic("unimplemented")
}
