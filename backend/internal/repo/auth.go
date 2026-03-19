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
	GetUsers(ctx context.Context, userIDs []string) ([]domain.User, error)
	GetUsersByTeamID(ctx context.Context, teamID string) ([]domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID string, hashedPassword string) error
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

func (i *userRepo) GetUsers(ctx context.Context, userIDs []string) ([]domain.User, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

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
			WHERE u.id = ANY(@user_ids);
			`

	rows, err := i.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"user_ids": userIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("collect users error: %w", err)
	}

	return users, nil
}
 
func (i *userRepo) GetUsersByTeamID(ctx context.Context, teamID string) ([]domain.User, error) {
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
			WHERE u.team_id = @team_id;
			`
 
	rows, err := i.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"team_id": teamID,
	})
	if err != nil {
		return nil, fmt.Errorf("query team users: %w", err)
	}
 
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("collect team users error: %w", err)
	}
 
	return users, nil
}

func (i *userRepo) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
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
			WHERE u.id = @id;
			`

	rows, err := i.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"id": id,
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

func (i *userRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE auth.t_users
		SET
			first_name = @first_name,
			last_name = @last_name,
			email = @email,
			updated_at = NOW()
		WHERE id = @id;
	`

	_, err := i.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	})
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (i *userRepo) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `
		UPDATE auth.t_users
		SET
			password_hash = @password_hash,
			updated_at = NOW()
		WHERE id = @id;
	`

	_, err := i.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"id":            userID,
		"password_hash": hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}

func NewUserRepo(dbClient *db.PostgresClient) UserRepository {
	return &userRepo{dbClient: dbClient}
}
