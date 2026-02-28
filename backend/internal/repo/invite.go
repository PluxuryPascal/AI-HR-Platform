package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrInviteNotFound = errors.New("invite not found")

type InviteRepository interface {
	CreateInvite(ctx context.Context, invite *domain.Invite, jobIDs []string) error
	GetInviteByToken(ctx context.Context, token string) (*domain.Invite, error)
	AcceptInviteAndCreateUser(ctx context.Context, inviteID string, user *domain.CreateUserRepoParams) (*domain.User, error)
}

type inviteRepo struct {
	dbClient *db.PostgresClient
}

func NewInviteRepo(dbClient *db.PostgresClient) InviteRepository {
	return &inviteRepo{dbClient: dbClient}
}

func (r *inviteRepo) CreateInvite(ctx context.Context, invite *domain.Invite, jobIDs []string) error {
	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const insertInvite = `
		INSERT INTO auth.t_invites (team_id, email, role, token, expires_at)
		VALUES (@team_id, @email, @role, @token, @expires_at)
		RETURNING id
	`

	if err := tx.QueryRow(ctx, insertInvite, pgx.NamedArgs{
		"team_id":    invite.TeamID,
		"email":      invite.Email,
		"role":       invite.Role,
		"token":      invite.Token,
		"expires_at": invite.ExpiresAt,
	}).Scan(&invite.ID); err != nil {
		return fmt.Errorf("insert invite: %w", err)
	}

	if len(jobIDs) > 0 {
		rows := make([][]any, len(jobIDs))
		for i, jobID := range jobIDs {
			rows[i] = []any{invite.ID, jobID}
		}

		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"auth", "t_invite_job_access"},
			[]string{"invite_id", "job_id"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return fmt.Errorf("batch insert invite job access: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *inviteRepo) GetInviteByToken(ctx context.Context, token string) (*domain.Invite, error) {
	const query = `
		SELECT id, team_id, email, role, token, expires_at, created_at
		FROM auth.t_invites
		WHERE token = @token
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"token": token})
	if err != nil {
		return nil, fmt.Errorf("query invite: %w", err)
	}

	invite, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Invite])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInviteNotFound
		}

		return nil, fmt.Errorf("scan invite: %w", err)
	}

	return &invite, nil
}

func (r *inviteRepo) AcceptInviteAndCreateUser(
	ctx context.Context,
	inviteID string,
	user *domain.CreateUserRepoParams,
) (*domain.User, error) {
	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const insertUser = `
		INSERT INTO auth.t_users (team_id, email, first_name, last_name, role, password_hash)
		VALUES (@team_id, @email, @first_name, @last_name, @role, @password_hash)
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
			COALESCE(locale, '') AS locale
	`

	userRows, err := tx.Query(ctx, insertUser, pgx.NamedArgs{
		"team_id":       user.TeamID,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"role":          user.Role,
		"password_hash": user.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	createdUser, err := pgx.CollectExactlyOneRow(userRows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	const transferAccess = `
		INSERT INTO hiring.t_job_access (user_id, job_id)
		SELECT @user_id, job_id
		FROM auth.t_invite_job_access
		WHERE invite_id = @invite_id
	`

	if _, err := tx.Exec(ctx, transferAccess, pgx.NamedArgs{
		"user_id":   createdUser.ID,
		"invite_id": inviteID,
	}); err != nil {
		return nil, fmt.Errorf("transfer job access: %w", err)
	}

	const deleteInvite = `DELETE FROM auth.t_invites WHERE id = @id`

	if _, err := tx.Exec(ctx, deleteInvite, pgx.NamedArgs{"id": inviteID}); err != nil {
		return nil, fmt.Errorf("delete invite: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &createdUser, nil
}
