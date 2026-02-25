package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Login(ctx context.Context, login string) (*domain.UserLogin, error)
	Register(ctx context.Context, user *domain.UserRegister) (*domain.User, error)
}

var ErrUserNotFound = errors.New("user not found")

type userRepo struct {
	dbClient *db.PostgresClient
}

func NewUserRepo(dbClient *db.PostgresClient) UserRepository {
	return &userRepo{dbClient: dbClient}
}

func (r *userRepo) Login(ctx context.Context, login string) (*domain.UserLogin, error) {
	query := `
		SELECT
			u.id AS user_id,
			u.password_hash AS password,
			t.name AS goup_alias
		FROM auth.t_users u
			JOIN auth.t_teams t on t.id = u.team_id
		WHERE u.email = @login;
	`

	result, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"login": login,
	})
	if err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}

	user, err := pgx.RowToStructByName[domain.UserLogin](result)
	if err != nil {
		return nil, fmt.Errorf("scan row error: %w", err)
	}

	return &user, nil
}

func (r *userRepo) Register(ctx context.Context, user *domain.UserRegister) (*domain.User, error) {
	query := `
		INSERT INTO auth.t.users u (
			email,
			password_hash
		)
		VALUES (
			@email,
			@password_hash
		)
	`

	result, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"email":         user.Mail,
		"password_hash": user.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}

	createdUser, err := pgx.RowToStructByName[domain.User](result)
	if err != nil {
		return nil, fmt.Errorf("scan row error: %w", err)
	}

	return &createdUser, nil
}
