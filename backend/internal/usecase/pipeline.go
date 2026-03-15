package usecase

import (
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
	UpdateStage(ctx context.Context, stage *domain.PipelineStage) error
	DeleteStage(ctx context.Context, id string) error
}

type pipelineUseCase struct {
	pipelineRepo repo.PipelineRepository
}

func NewPipelineUseCase(pipelineRepo repo.PipelineRepository) PipelineUseCase {
	return &pipelineUseCase{pipelineRepo: pipelineRepo}
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

func (u *pipelineUseCase) UpdateStage(ctx context.Context, stage *domain.PipelineStage) error {
	if err := u.pipelineRepo.UpdateStage(ctx, stage); err != nil {
		return fmt.Errorf("update pipeline stage: %w", err)
	}
	return nil
}

func (u *pipelineUseCase) DeleteStage(ctx context.Context, id string) error {
	if err := u.pipelineRepo.DeleteStage(ctx, id); err != nil {
		return fmt.Errorf("delete pipeline stage: %w", err)
	}
	return nil
}
