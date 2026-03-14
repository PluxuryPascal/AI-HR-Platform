package usecase

import (
	"backend/internal/domain"
	"backend/internal/repo"
	"context"
	"errors"
	"fmt"
)

var ErrJobNotFound = errors.New("job not found")

type JobUseCase interface {
	CreateJob(ctx context.Context, job *domain.Job) error
	GetJobByID(ctx context.Context, id string) (*domain.Job, error)
	GetJobsByTeam(ctx context.Context, teamID string, offset, limit int, filter domain.JobFilter) (*domain.JobsDTO, error)
	UpdateJob(ctx context.Context, job *domain.Job) error
	DeleteJob(ctx context.Context, id string) error
}

type jobUseCase struct {
	jobRepo repo.JobRepository
}

func NewJobUseCase(jobRepo repo.JobRepository) JobUseCase {
	return &jobUseCase{jobRepo: jobRepo}
}

func (u *jobUseCase) CreateJob(ctx context.Context, job *domain.Job) error {
	if err := u.jobRepo.Create(ctx, job); err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	return nil
}

func (u *jobUseCase) GetJobByID(ctx context.Context, id string) (*domain.Job, error) {
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, fmt.Errorf("get job by id: %w", err)
	}

	return job, nil
}

func (u *jobUseCase) GetJobsByTeam(ctx context.Context, teamID string, offset, limit int, filter domain.JobFilter) (*domain.JobsDTO, error) {
	jobsDTO, err := u.jobRepo.GetByTeamID(ctx, teamID, offset, limit, filter)
	if err != nil {
		return nil, fmt.Errorf("get jobs by team: %w", err)
	}

	return jobsDTO, nil
}

func (u *jobUseCase) UpdateJob(ctx context.Context, job *domain.Job) error {
	if err := u.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("update job: %w", err)
	}

	return nil
}

func (u *jobUseCase) DeleteJob(ctx context.Context, id string) error {
	if err := u.jobRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete job: %w", err)
	}

	return nil
}
