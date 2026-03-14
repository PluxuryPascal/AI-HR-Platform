package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ProcessingRepository interface {
	CreateTask(ctx context.Context, task *domain.ProcessingTask) error
	UpdateStatus(ctx context.Context, id string, status domain.TaskStatus, progress int, errMsg *string) error
}

type processingRepo struct {
	dbClient *db.PostgresClient
}

func NewProcessingRepo(dbClient *db.PostgresClient) ProcessingRepository {
	return &processingRepo{dbClient: dbClient}
}

func (r *processingRepo) CreateTask(ctx context.Context, task *domain.ProcessingTask) error {
	const query = `
		INSERT INTO ai_engine.t_processing_tasks (workflow_id, entity_id, status, progress_percent)
		VALUES (@workflow_id, @entity_id, @status, @progress)
		RETURNING id, updated_at
	`
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"workflow_id": task.WorkflowID,
		"entity_id":   task.EntityID,
		"status":      task.Status,
		"progress":    task.ProgressPercent,
	}).Scan(&task.ID, &task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create processing task: %w", err)
	}
	return nil
}

func (r *processingRepo) UpdateStatus(ctx context.Context, id string, status domain.TaskStatus, progress int, errMsg *string) error {
	const query = `
		UPDATE ai_engine.t_processing_tasks
		SET status = @status, progress_percent = @progress, error_message = @error, updated_at = NOW()
		WHERE id = @id
	`
	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"id":       id,
		"status":   status,
		"progress": progress,
		"error":    errMsg,
	})
	if err != nil {
		return fmt.Errorf("update processing status: %w", err)
	}
	return nil
}
