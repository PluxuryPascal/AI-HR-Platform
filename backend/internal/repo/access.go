package repo

import (
	"backend/internal/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type AccessRepository interface {
	TransferAccess(ctx context.Context, userID string, jobIDs []string) error
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

func NewAccessRepository(dbClient *db.PostgresClient) AccessRepository {
	return &accessRepository{dbClient: dbClient}
}
