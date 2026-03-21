package usecase

import (
	"backend/internal/audit"
	"backend/internal/domain"
	"backend/internal/repo"
	"context"
	"errors"
	"fmt"
)

var ErrJobNotFound = errors.New("job not found")

type JobUseCase interface {
	CreateJob(ctx context.Context, job *domain.Job) error
	GetJobByID(ctx context.Context, id, userID, role string) (*domain.Job, error)
	GetJobsByTeam(ctx context.Context, teamID, userID, role string, offset, limit int, filter domain.JobFilter) (*domain.JobsDTO, error)
	UpdateJob(ctx context.Context, job *domain.Job, actorID, role string) error
	DeleteJob(ctx context.Context, id, actorID, teamID, role string) error
}

type jobUseCase struct {
	jobRepo    repo.JobRepository
	accessRepo repo.AccessRepository
	auditor    *audit.Logger
}

func NewJobUseCase(jobRepo repo.JobRepository, accessRepo repo.AccessRepository, auditor *audit.Logger) JobUseCase {
	return &jobUseCase{
		jobRepo:    jobRepo,
		accessRepo: accessRepo,
		auditor:    auditor,
	}
}

func (u *jobUseCase) CreateJob(ctx context.Context, job *domain.Job) error {
	if err := u.jobRepo.Create(ctx, job); err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    job.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &job.CreatedBy,
		Action:    audit.HiringJobCreated,
		TargetID:  &job.ID,
	})

	return nil
}

func (u *jobUseCase) GetJobByID(ctx context.Context, id, userID, role string) (*domain.Job, error) {
	job, err := u.jobRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, fmt.Errorf("get job by id: %w", err)
	}

	if role != "owner" && role != "admin" {
		ok, err := u.accessRepo.HasAccess(ctx, userID, id)
		if err != nil {
			return nil, fmt.Errorf("check access: %w", err)
		}
		if !ok {
			return nil, errors.New("access denied to this job")
		}
	}

	return job, nil
}

func (u *jobUseCase) GetJobsByTeam(ctx context.Context, teamID, userID, role string, offset, limit int, filter domain.JobFilter) (*domain.JobsDTO, error) {
	if role != "owner" && role != "admin" {
		filter.AllowedUserID = &userID
	}

	jobsDTO, err := u.jobRepo.GetByTeamID(ctx, teamID, offset, limit, filter)
	if err != nil {
		return nil, fmt.Errorf("get jobs by team: %w", err)
	}

	return jobsDTO, nil
}

func (u *jobUseCase) UpdateJob(ctx context.Context, job *domain.Job, actorID, role string) error {
	// Let's get the old job to check status transition for logging
	oldJob, err := u.GetJobByID(ctx, job.ID, actorID, role)
	if err != nil {
		return fmt.Errorf("get old job: %w", err)
	}

	if err := u.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("update job: %w", err)
	}

	if oldJob.Status != job.Status {
		statusAction := map[domain.JobStatus]string{
			domain.JobStatusPublished: audit.HiringJobPublished,
			domain.JobStatusClosed:    audit.HiringJobClosed,
			domain.JobStatusArchived:  audit.HiringJobArchived,
		}
		if action, ok := statusAction[job.Status]; ok {
			_ = u.auditor.Log(ctx, audit.Entry{
				TeamID:    job.TeamID,
				ActorType: audit.ActorUser,
				ActorID:   &actorID,
				Action:    action,
				TargetID:  &job.ID,
			})
		}
	}

	return nil
}

func (u *jobUseCase) DeleteJob(ctx context.Context, id, actorID, teamID, role string) error {
	if _, err := u.GetJobByID(ctx, id, actorID, role); err != nil {
		return err
	}

	if err := u.jobRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete job: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    teamID,
		ActorType: audit.ActorUser,
		ActorID:   &actorID,
		Action:    audit.HiringJobDeleted,
		TargetID:  &id,
	})

	return nil
}
