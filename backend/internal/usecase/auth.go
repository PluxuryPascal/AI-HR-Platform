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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (*string, time.Duration, error)
	RegisterOwner(ctx context.Context, req domain.RegisterOwnerRequest) (*string, time.Duration, error)
	Logout(ctx context.Context, tokenStr string) error
}

type authUseCase struct {
	repo         repo.UserRepository
	cacheManager *cache.Manager
	token        *token.JWTtoken
	hash         hash.Hash
}

func (a *authUseCase) Logout(ctx context.Context, tokenStr string) error {
	token, err := a.token.VerifyToken([]byte(tokenStr))
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	sessionID := token.Subject()

	if err := cache.Delete(ctx, a.cacheManager, cache.SessionKey, sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (a *authUseCase) RegisterOwner(ctx context.Context, req domain.RegisterOwnerRequest) (*string, time.Duration, error) {
	hashedPassword, err := a.hash.Hash(req.Password)
	if err != nil {
		return nil, 0, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.RegisterOwnerRequest{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		TeamName:  req.TeamName,
	}

	userData, err := a.repo.RegisterOwner(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, 0, ErrUserAlreadyExists
		}

		return nil, 0, fmt.Errorf("repo register: %w", err)
	}

	sessionID := uuid.New().String()

	if err := cache.SetWithTTL(ctx, a.cacheManager, cache.SessionKey, sessionID, domain.Session{
		UserID: userData.ID,
		TeamID: userData.TeamID,
		Role:   userData.Role,
	}, a.token.ExpireAt); err != nil {
		return nil, 0, fmt.Errorf("set session: %w", err)
	}

	signed, err := a.token.GenerateToken(sessionID)
	if err != nil {
		return nil, 0, fmt.Errorf("generate token: %w", err)
	}

	tokenString := string(signed)

	return &tokenString, a.token.ExpireAt, nil
}

func (a *authUseCase) Login(ctx context.Context, email string, password string) (*string, time.Duration, error) {
	user, err := a.repo.Login(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, 0, ErrInvalidCredentials
		}

		return nil, 0, fmt.Errorf("repo login: %w", err)
	}

	verified, err := a.hash.Verify(password, user.PasswordHash)
	if err != nil {
		return nil, 0, fmt.Errorf("verify password hash: %w", err)
	}

	if !verified {
		return nil, 0, ErrInvalidCredentials
	}

	sessionID := uuid.New().String()

	if err := cache.SetWithTTL(ctx, a.cacheManager, cache.SessionKey, sessionID, domain.Session{
		UserID: user.ID,
		TeamID: user.TeamID,
		Role:   user.Role,
	}, a.token.ExpireAt); err != nil {
		return nil, 0, fmt.Errorf("set session: %w", err)
	}

	signed, err := a.token.GenerateToken(sessionID)
	if err != nil {
		return nil, 0, fmt.Errorf("generate token: %w", err)
	}

	tokenString := string(signed)

	return &tokenString, a.token.ExpireAt, nil
}

func NewAuthUseCase(repo repo.UserRepository, cacheManager *cache.Manager, token *token.JWTtoken, hash hash.Hash) AuthUseCase {
	return &authUseCase{
		repo:         repo,
		cacheManager: cacheManager,
		token:        token,
		hash:         hash,
	}
}

var _ AuthUseCase = (*authUseCase)(nil)
