package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PipelineRepository interface {
	CreateStage(ctx context.Context, params domain.CreateStageParams) (*domain.PipelineStage, error)
	GetStagesByJobID(ctx context.Context, jobID string) ([]domain.PipelineStage, error)
	GetStagesByTeamID(ctx context.Context, teamID string) ([]domain.PipelineStage, error)
	UpdateStage(ctx context.Context, stage *domain.PipelineStage) error
	DeleteStage(ctx context.Context, id string) error
	GetStageByID(ctx context.Context, id string) (*domain.PipelineStage, error)
}

type pipelineRepo struct {
	dbClient *db.PostgresClient
}

func NewPipelineRepo(dbClient *db.PostgresClient) PipelineRepository {
	return &pipelineRepo{dbClient: dbClient}
}

func (r *pipelineRepo) CreateStage(ctx context.Context, params domain.CreateStageParams) (*domain.PipelineStage, error) {
	const query = `
		INSERT INTO hiring.t_pipeline_stages (job_id, team_id, code, title, position, is_terminal, is_rejection, is_interview, color)
		VALUES (@job_id, @team_id, @code, @title, @position, @is_terminal, @is_rejection, @is_interview, @color)
		RETURNING id, job_id, team_id, code, title, position, is_terminal, is_rejection, is_interview, color
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"job_id":       params.JobID,
		"team_id":      params.TeamID,
		"code":         params.Code,
		"title":        params.Title,
		"position":     params.Position,
		"is_terminal":  params.IsTerminal,
		"is_rejection": params.IsRejection,
		"is_interview": params.IsInterview,
		"color":        params.Color,
	})
	if err != nil {
		return nil, fmt.Errorf("insert stage: %w", err)
	}

	stage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.PipelineStage])
	if err != nil {
		return nil, fmt.Errorf("scan stage: %w", err)
	}

	return &stage, nil
}

func (r *pipelineRepo) GetStagesByJobID(ctx context.Context, jobID string) ([]domain.PipelineStage, error) {
	const query = `
		SELECT id, job_id, team_id, code, title, position, is_terminal, is_rejection, is_interview, color
		FROM hiring.t_pipeline_stages
		WHERE job_id = @job_id
		ORDER BY position ASC
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"job_id": jobID})
	if err != nil {
		return nil, fmt.Errorf("query stages by job: %w", err)
	}

	stages, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.PipelineStage])
	if err != nil {
		return nil, fmt.Errorf("scan stages: %w", err)
	}

	return stages, nil
}

func (r *pipelineRepo) GetStagesByTeamID(ctx context.Context, teamID string) ([]domain.PipelineStage, error) {
	const query = `
		SELECT id, job_id, team_id, code, title, position, is_terminal, is_rejection, is_interview, color
		FROM hiring.t_pipeline_stages
		WHERE team_id = @team_id AND job_id IS NULL
		ORDER BY position ASC
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"team_id": teamID})
	if err != nil {
		return nil, fmt.Errorf("query stages by team: %w", err)
	}

	stages, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.PipelineStage])
	if err != nil {
		return nil, fmt.Errorf("scan stages: %w", err)
	}

	return stages, nil
}

func (r *pipelineRepo) GetStageByID(ctx context.Context, id string) (*domain.PipelineStage, error) {
	const query = `
		SELECT id, job_id, team_id, code, title, position, is_terminal, is_rejection, is_interview, color
		FROM hiring.t_pipeline_stages
		WHERE id = @id
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("query stage by id: %w", err)
	}

	stage, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.PipelineStage])
	if err != nil {
		return nil, fmt.Errorf("scan stage: %w", err)
	}

	return &stage, nil
}

func (r *pipelineRepo) UpdateStage(ctx context.Context, stage *domain.PipelineStage) error {
	const query = `
		UPDATE hiring.t_pipeline_stages
		SET title = @title, position = @position, is_terminal = @is_terminal, is_rejection = @is_rejection, is_interview = @is_interview, color = @color
		WHERE id = @id
	`

	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"id":           stage.ID,
		"title":        stage.Title,
		"position":     stage.Position,
		"is_terminal":  stage.IsTerminal,
		"is_rejection": stage.IsRejection,
		"is_interview": stage.IsInterview,
		"color":        stage.Color,
	})
	if err != nil {
		return fmt.Errorf("update stage: %w", err)
	}

	return nil
}

func (r *pipelineRepo) DeleteStage(ctx context.Context, id string) error {
	const query = `DELETE FROM hiring.t_pipeline_stages WHERE id = @id`
	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("delete stage: %w", err)
	}
	return nil
}
