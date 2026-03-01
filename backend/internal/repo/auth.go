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
	Login(ctx context.Context, login string) (*domain.User, error)
	RegisterOwner(ctx context.Context, user *domain.RegisterOwnerRequest) (*domain.User, error)
}

var ErrUserNotFound = errors.New("user not found")

type userRepo struct {
	dbClient *db.PostgresClient
}

func (i *userRepo) Login(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT
			u.id,
			t.id AS team_id,
			t.name AS team_name,
			u.email,
			u.first_name,
			u.last_name,
			u.role,
			u.password_hash,
			u.created_at,
			u.updated_at,
			COALESCE(u.locale, '') AS locale
			FROM auth.t_users u
			JOIN auth.t_teams t on t.id = u.team_id
			WHERE u.email = @email;
			`

	rows, err := i.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("scan row error: %w", err)
	}

	return &user, nil
}

func (i *userRepo) RegisterOwner(ctx context.Context, user *domain.RegisterOwnerRequest) (*domain.User, error) {
	tx, err := i.dbClient.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx error: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var teamID string

	teamQuery := `
		INSERT INTO auth.t_teams (
			name
		)
		VALUES (
			@name
		)
		RETURNING id;
	`

	if err := tx.QueryRow(ctx, teamQuery, pgx.NamedArgs{
		"name": user.TeamName,
	}).Scan(&teamID); err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}

	query := `
		INSERT INTO auth.t_users (
			team_id,
			email,
			first_name,
			last_name,
			role,
			password_hash
		)
		VALUES (
			@team_id,
			@email,
			@first_name,
			@last_name,
			'owner',
			@password_hash
		)
		RETURNING 
			id,
			(SELECT t.id   FROM auth.t_teams t WHERE t.id = team_id) AS team_id,
			(SELECT t.name FROM auth.t_teams t WHERE t.id = team_id) AS team_name,
			email,
			first_name,
			last_name,
			role,
			password_hash,
			created_at,
			updated_at,
			COALESCE(locale, '') AS locale;
	`

	rows, err := tx.Query(ctx, query, pgx.NamedArgs{
		"team_id":       teamID,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"password_hash": user.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}

	createdUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("scan row error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx error: %w", err)
	}

	return &createdUser, nil
}

func NewUserRepo(dbClient *db.PostgresClient) UserRepository {
	return &userRepo{dbClient: dbClient}
}
