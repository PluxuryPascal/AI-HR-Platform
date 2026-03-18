package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type AccessRepository interface {
	TransferAccess(ctx context.Context, userID string, jobIDs []string) error
	GrantAccess(ctx context.Context, userID, jobID string) error
	RevokeAccess(ctx context.Context, userID, jobID string) error
	GetAccessByJobID(ctx context.Context, jobID string) ([]domain.JobAccess, error)
}

type accessRepository struct {
	dbClient *db.PostgresClient
}

func (r *accessRepository) TransferAccess(ctx context.Context, userID string, jobIDs []string) error {
	if len(jobIDs) == 0 {
		return nil
	}

	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows := make([][]any, len(jobIDs))
	for i, jobID := range jobIDs {
		rows[i] = []any{userID, jobID}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"hiring", "t_job_access"},
		[]string{"user_id", "job_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("batch insert job access: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *accessRepository) GrantAccess(ctx context.Context, userID, jobID string) error {
	const query = `
		INSERT INTO hiring.t_job_access (user_id, job_id)
		VALUES (@user_id, @job_id)
		ON CONFLICT (user_id, job_id) DO NOTHING
	`
	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"user_id": userID,
		"job_id":  jobID,
	})
	if err != nil {
		return fmt.Errorf("grant access: %w", err)
	}
	return nil
}

func (r *accessRepository) RevokeAccess(ctx context.Context, userID, jobID string) error {
	const query = `DELETE FROM hiring.t_job_access WHERE user_id = @user_id AND job_id = @job_id`
	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"user_id": userID,
		"job_id":  jobID,
	})
	if err != nil {
		return fmt.Errorf("revoke access: %w", err)
	}
	return nil
}

func (r *accessRepository) GetAccessByJobID(ctx context.Context, jobID string) ([]domain.JobAccess, error) {
	const query = `
		SELECT user_id, job_id
		FROM hiring.t_job_access
		WHERE job_id = @job_id
	`
	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"job_id": jobID})
	if err != nil {
		return nil, fmt.Errorf("query job access: %w", err)
	}
	defer rows.Close()

	var list []domain.JobAccess
	for rows.Next() {
		var a domain.JobAccess
		if err := rows.Scan(&a.UserID, &a.JobID); err != nil {
			return nil, fmt.Errorf("scan job access: %w", err)
		}
		list = append(list, a)
	}
	return list, nil
}

func NewAccessRepository(dbClient *db.PostgresClient) AccessRepository {
	return &accessRepository{dbClient: dbClient}
}

