package usecase

import (
	"backend/internal/audit"
	"backend/internal/cache"
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/pkg/hash"
	"backend/pkg/rbac"
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
	GetTeamMembers(ctx context.Context, teamID string) ([]domain.User, error)
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, firstName, lastName, email string) error
	UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error
	DeleteMember(ctx context.Context, teamID, memberID, actorID, actorRole string) error
}

type authUseCase struct {
	repo         repo.UserRepository
	cacheManager *cache.Manager
	token        *token.JWTtoken
	hash         hash.Hash
	enforcer     *rbac.CasbinClient
	auditor      *audit.Logger
}

func (a *authUseCase) Logout(ctx context.Context, tokenStr string) error {
	token, err := a.token.VerifyToken([]byte(tokenStr))
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	sessionID := token.Subject()

	sess, err := cache.Get(ctx, a.cacheManager, cache.SessionKey, sessionID)
	if err == nil {
		_ = a.auditor.Log(ctx, audit.Entry{
			TeamID:    sess.TeamID,
			ActorType: audit.ActorUser,
			ActorID:   &sess.UserID,
			Action:    audit.AuthUserLoggedOut,
		})
	}

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

	if _, err := a.enforcer.AddRoleForUserInDomain(userData.ID, "owner", userData.TeamID); err != nil {
		return nil, 0, fmt.Errorf("add grouping policy: %w", err)
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

	_ = a.auditor.Log(ctx, audit.Entry{
		TeamID:    userData.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &userData.ID,
		Action:    audit.AuthOwnerRegistered,
		TargetID:  &userData.TeamID,
	})

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

	_ = a.auditor.Log(ctx, audit.Entry{
		TeamID:    user.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &user.ID,
		Action:    audit.AuthUserLoggedIn,
	})

	return &tokenString, a.token.ExpireAt, nil
}

func (a *authUseCase) GetTeamMembers(ctx context.Context, teamID string) ([]domain.User, error) {
	users, err := a.repo.GetUsersByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("repo get team users: %w", err)
	}
	return users, nil
}

func (a *authUseCase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := a.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repo get user: %w", err)
	}
	return user, nil
}

func (a *authUseCase) UpdateProfile(ctx context.Context, userID string, firstName, lastName, email string) error {
	user, err := a.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("repo get user: %w", err)
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email

	if err := a.repo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("repo update user: %w", err)
	}

	return nil
}

func (a *authUseCase) UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error {
	user, err := a.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("repo get user: %w", err)
	}

	verified, err := a.hash.Verify(currentPassword, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("verify current password: %w", err)
	}

	if !verified {
		return errors.New("invalid current password")
	}

	hashedPassword, err := a.hash.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	if err := a.repo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("repo update password: %w", err)
	}

	return nil
}

func (a *authUseCase) DeleteMember(ctx context.Context, teamID, memberID, actorID, actorRole string) error {
	if memberID == actorID {
		return errors.New("cannot delete yourself")
	}

	// Check if member belongs to the same team
	member, err := a.repo.GetUserByID(ctx, memberID)
	if err != nil {
		return fmt.Errorf("get member: %w", err)
	}

	if member.TeamID != teamID {
		return errors.New("member does not belong to your team")
	}

	// Check permissions: only owners/admins can delete
	if actorRole != "owner" && actorRole != "admin" {
		return errors.New("only owners or admins can delete team members")
	}

	// If deleting an owner, only another owner can do it (or maybe not at all, but let's allow owner to delete owner if not self)
	if member.Role == "owner" && actorRole != "owner" {
		return errors.New("only owners can delete other owners")
	}

	if err := a.repo.DeleteUser(ctx, memberID); err != nil {
		return fmt.Errorf("repo delete user: %w", err)
	}

	_ = a.auditor.Log(ctx, audit.Entry{
		TeamID:    teamID,
		ActorType: audit.ActorUser,
		ActorID:   &actorID,
		Action:    audit.HiringStageDeleted, // Reuse or add new audit action
		TargetID:  &memberID,
	})

	return nil
}

func NewAuthUseCase(repo repo.UserRepository, cacheManager *cache.Manager, token *token.JWTtoken, hash hash.Hash, enforcer *rbac.CasbinClient, auditor *audit.Logger) AuthUseCase {
	return &authUseCase{
		repo:         repo,
		cacheManager: cacheManager,
		token:        token,
		hash:         hash,
		enforcer:     enforcer,
		auditor:      auditor,
	}
}

var _ AuthUseCase = (*authUseCase)(nil)
