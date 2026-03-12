package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrInviteNotFound = errors.New("invite not found")

type InviteRepository interface {
	CreateInvite(ctx context.Context, invite *domain.Invite, jobIDs []string) error
	GetInviteByToken(ctx context.Context, token string) (*domain.Invite, error)
	LockInviteForProcessing(ctx context.Context, token string) (*domain.Invite, []string, error)
	AcceptInviteLocalTx(ctx context.Context, user *domain.CreateUserRepoParams) (*domain.User, error)
	UpdateInviteStatus(ctx context.Context, id string, status string) error
	RollbackInviteToPending(ctx context.Context, id string) error
	GetStuckInvites(ctx context.Context, threshold time.Duration) ([]domain.Invite, error)
	GetInviteJobIDs(ctx context.Context, inviteID string) ([]string, error)
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
		INSERT INTO auth.t_invites (team_id, email, role, token, expires_at, status)
		VALUES (@team_id, @email, @role, @token, @expires_at, 'pending')
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
		WHERE token = @token AND status = 'pending'
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

func (r *inviteRepo) LockInviteForProcessing(ctx context.Context, token string) (*domain.Invite, []string, error) {
	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const updateQuery = `
		UPDATE auth.t_invites 
		SET status = 'processing', updated_at = NOW() 
		WHERE token = @token AND status = 'pending' 
		RETURNING id, team_id, email, role, token, status, expires_at, created_at, updated_at
	`

	rows, err := tx.Query(ctx, updateQuery, pgx.NamedArgs{"token": token})
	if err != nil {
		return nil, nil, fmt.Errorf("update invite status: %w", err)
	}

	invite, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Invite])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrInviteNotFound
		}

		return nil, nil, fmt.Errorf("scan locked invite: %w", err)
	}

	const jobsQuery = `
		SELECT job_id
		FROM auth.t_invite_job_access
		WHERE invite_id = @invite_id
	`
	jobRows, err := tx.Query(ctx, jobsQuery, pgx.NamedArgs{"invite_id": invite.ID})
	if err != nil {
		return nil, nil, fmt.Errorf("query associated jobs: %w", err)
	}

	jobIDs, err := pgx.CollectRows(jobRows, pgx.RowTo[string])
	if err != nil {
		return nil, nil, fmt.Errorf("scan jobs: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("commit tx: %w", err)
	}

	return &invite, jobIDs, nil
}

func (r *inviteRepo) AcceptInviteLocalTx(ctx context.Context, user *domain.CreateUserRepoParams) (*domain.User, error) {
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

	rows, err := r.dbClient.Pool.Query(ctx, insertUser, pgx.NamedArgs{
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

	createdUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return &createdUser, nil
}

func (r *inviteRepo) UpdateInviteStatus(ctx context.Context, id string, status string) error {
	const updateQuery = `
		UPDATE auth.t_invites 
		SET status = @status, updated_at = NOW() 
		WHERE id = @id
	`

	if _, err := r.dbClient.Pool.Exec(ctx, updateQuery, pgx.NamedArgs{
		"status": status,
		"id":     id,
	}); err != nil {
		return fmt.Errorf("update invite status: %w", err)
	}

	return nil
}

func (r *inviteRepo) RollbackInviteToPending(ctx context.Context, id string) error {
	const query = `
		UPDATE auth.t_invites 
		SET status = 'pending', updated_at = NOW() 
		WHERE id = @id
	`
	if _, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{"id": id}); err != nil {
		return fmt.Errorf("rollback invite to pending: %w", err)
	}
	return nil
}
func (r *inviteRepo) GetStuckInvites(ctx context.Context, threshold time.Duration) ([]domain.Invite, error) {
	const query = `
		UPDATE auth.t_invites
		SET updated_at = NOW()
		WHERE status = 'processing' AND updated_at < @threshold_time
		RETURNING id, team_id, email, role, token, status, expires_at, created_at, updated_at
	`

	thresholdTime := time.Now().Add(-threshold)
	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"threshold_time": thresholdTime})
	if err != nil {
		return nil, fmt.Errorf("query stuck invites: %w", err)
	}

	invites, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Invite])
	if err != nil {
		return nil, fmt.Errorf("collect stuck invites: %w", err)
	}

	return invites, nil
}

func (r *inviteRepo) GetInviteJobIDs(ctx context.Context, inviteID string) ([]string, error) {
	const query = `
		SELECT job_id
		FROM auth.t_invite_job_access
		WHERE invite_id = @invite_id
	`
	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"invite_id": inviteID})
	if err != nil {
		return nil, fmt.Errorf("query job ids: %w", err)
	}

	jobIDs, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("collect job ids: %w", err)
	}

	return jobIDs, nil
}
