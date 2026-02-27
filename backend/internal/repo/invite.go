package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type InviteRepository interface {
	CreateUser(ctx context.Context, user *domain.CreateUserRepoParams) (*domain.User, error)
}

type inviteRepo struct {
	dbClient *db.PostgresClient
}

func NewInviteRepo(dbClient *db.PostgresClient) InviteRepository {
	return &inviteRepo{dbClient: dbClient}
}

func (i *inviteRepo) CreateUser(ctx context.Context, user *domain.CreateUserRepoParams) (*domain.User, error) {
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
			@role,
			@password_hash
		)
		RETURNING 
			u.id,
			(SELECT t.id FROM auth.t_teams t WHERE t.id = team_id) AS team_id,
			(SELECT t.name FROM auth.t_teams t WHERE t.id = team_id) AS team_name,
			u.email,
			u.first_name,
			u.last_name,
			u.role,
			u.password_hash,
			u.created_at,
			u.updated_at,
			u.locale;
	`

	rows, err := i.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"team_id":       user.TeamID,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"role":          user.Role,
		"password_hash": user.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}

	createdUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("scan row error: %w", err)
	}

	return &createdUser, nil
}
