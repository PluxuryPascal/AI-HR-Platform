package usecase

import (
	"backend/internal/audit"
	"backend/internal/domain"
	"backend/internal/repo"
	"context"
	"errors"
	"fmt"
)

var ErrNoPipelineStages = errors.New("no pipeline stages")

type PipelineUseCase interface {
	CreateStage(ctx context.Context, params domain.CreateStageParams) (*domain.PipelineStage, error)
	GetStagesByJobID(ctx context.Context, jobID string) ([]domain.PipelineStage, error)
	GetStagesByTeamID(ctx context.Context, teamID string) ([]domain.PipelineStage, error)
	GetFirstStageByJobID(ctx context.Context, jobID string) (*domain.PipelineStage, error)
	UpdateStage(ctx context.Context, stage *domain.PipelineStage, actorID string) error
	DeleteStage(ctx context.Context, id, actorID, teamID string) error
}

type pipelineUseCase struct {
	pipelineRepo  repo.PipelineRepository
	candidateRepo repo.CandidateRepository
	auditor       *audit.Logger
}

func NewPipelineUseCase(pipelineRepo repo.PipelineRepository, candidateRepo repo.CandidateRepository, auditor *audit.Logger) PipelineUseCase {
	return &pipelineUseCase{
		pipelineRepo:  pipelineRepo,
		candidateRepo: candidateRepo,
		auditor:       auditor,
	}
}

func (u *pipelineUseCase) GetFirstStageByJobID(ctx context.Context, jobID string) (*domain.PipelineStage, error) {
	stages, err := u.pipelineRepo.GetStagesByJobID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get stages: %w", err)
	}
	if len(stages) == 0 {
		return nil, ErrNoPipelineStages
	}

	return &stages[0], nil
}

func (u *pipelineUseCase) CreateStage(ctx context.Context, params domain.CreateStageParams) (*domain.PipelineStage, error) {
	stage, err := u.pipelineRepo.CreateStage(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create pipeline stage: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    params.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &params.ActorID,
		Action:    audit.HiringStageCreated,
		TargetID:  &stage.ID,
	})

	return stage, nil
}

func (u *pipelineUseCase) GetStagesByJobID(ctx context.Context, jobID string) ([]domain.PipelineStage, error) {
	stages, err := u.pipelineRepo.GetStagesByJobID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get stages by job: %w", err)
	}
	return stages, nil
}

func (u *pipelineUseCase) GetStagesByTeamID(ctx context.Context, teamID string) ([]domain.PipelineStage, error) {
	stages, err := u.pipelineRepo.GetStagesByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("get stages by team: %w", err)
	}
	return stages, nil
}

func (u *pipelineUseCase) UpdateStage(ctx context.Context, stage *domain.PipelineStage, actorID string) error {
	if err := u.pipelineRepo.UpdateStage(ctx, stage); err != nil {
		return fmt.Errorf("update pipeline stage: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    stage.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &actorID,
		Action:    audit.HiringStageUpdated,
		TargetID:  &stage.ID,
	})

	return nil
}

func (u *pipelineUseCase) DeleteStage(ctx context.Context, id, actorID, teamID string) error {
	// 1. Get the stage to be deleted to find its JobID
	stage, err := u.pipelineRepo.GetStageByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get stage to delete: %w", err)
	}

	// 2. If it's a job-specific stage, we must move candidates
	if stage.JobID != nil {
		stages, err := u.pipelineRepo.GetStagesByJobID(ctx, *stage.JobID)
		if err != nil {
			return fmt.Errorf("get job stages: %w", err)
		}

		if len(stages) <= 1 {
			return fmt.Errorf("cannot delete the only stage in a pipeline")
		}

		// Pick the first stage as the target for moved candidates
		targetStageID := stages[0].ID
		if targetStageID == id {
			// If we are deleting the first stage, pick the second one
			targetStageID = stages[1].ID
		}

		// 3. Move candidates before deleting the stage
		if err := u.candidateRepo.MoveAllToStage(ctx, id, targetStageID); err != nil {
			return fmt.Errorf("move candidates before deletion: %w", err)
		}
	}

	// 4. Delete the stage
	if err := u.pipelineRepo.DeleteStage(ctx, id); err != nil {
		return fmt.Errorf("delete pipeline stage: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    teamID,
		ActorType: audit.ActorUser,
		ActorID:   &actorID,
		Action:    audit.HiringStageDeleted,
		TargetID:  &id,
	})

	return nil
}
