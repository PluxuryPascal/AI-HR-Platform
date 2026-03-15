package usecase

import (
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/pkg/mq"
	"backend/pkg/storage"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrCandidateNotFound  = errors.New("candidate not found")
	ErrStorageUnavailable = errors.New("storage unavailable")
)

type CandidateUseCase interface {
	CreateCandidate(ctx context.Context, params domain.CreateCandidateParams) (*domain.Candidate, error)
	GetCandidateByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error)
	GetCandidatesByJob(ctx context.Context, jobID string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error)
	UpdateCandidate(ctx context.Context, candidate *domain.Candidate) error
	UpdateCandidateProfile(ctx context.Context, profile *domain.CandidateProfile) error
	DeleteCandidate(ctx context.Context, id string) error
	MoveCandidate(ctx context.Context, p domain.MoveCandidateParams) error
}

type candidateUseCase struct {
	log           *zap.Logger
	candidateRepo repo.CandidateRepository
	pipelineRepo  repo.PipelineRepository
	storage       storage.FileStorage
	publisher     *mq.MQPublisher
}

func NewCandidateUseCase(log *zap.Logger, candidateRepo repo.CandidateRepository, pipelineRepo repo.PipelineRepository, storage storage.FileStorage, publisher *mq.MQPublisher) CandidateUseCase {
	return &candidateUseCase{
		log:           log,
		candidateRepo: candidateRepo,
		pipelineRepo:  pipelineRepo,
		storage:       storage,
		publisher:     publisher,
	}
}

func (u *candidateUseCase) CreateCandidate(ctx context.Context, params domain.CreateCandidateParams) (*domain.Candidate, error) {
	stages, err := u.pipelineRepo.GetStagesByJobID(ctx, params.JobID)
	if err != nil {
		return nil, fmt.Errorf("get stages: %w", err)
	}

	if len(stages) == 0 {
		return nil, ErrNoPipelineStages
	}

	initialStageID := stages[0].ID

	fileKey, err := u.storage.UploadFile(ctx, params.File, params.Filename)
	if err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	candidate := &domain.Candidate{
		JobID:         params.JobID,
		ResumeFileKey: &fileKey,
		ParsingStatus: domain.ParsingStatusPending,
	}

	candidate, err = u.candidateRepo.Create(ctx, candidate)
	if err != nil {
		return nil, fmt.Errorf("create candidate: %w", err)
	}

	return candidate, nil
}

func (u *candidateUseCase) GetCandidateByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error) {
	cand, profile, stageID, err := u.candidateRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, nil, ErrCandidateNotFound
		}

		return nil, nil, nil, fmt.Errorf("get candidate by id: %w", err)
	}

	return cand, profile, stageID, nil
}

func (u *candidateUseCase) GetCandidatesByJob(ctx context.Context, jobID string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error) {
	candidatesDTO, err := u.candidateRepo.GetByJobID(ctx, jobID, offset, limit, filter)
	if err != nil {
		return nil, fmt.Errorf("get candidates by job: %w", err)
	}

	return candidatesDTO, nil
}

func (u *candidateUseCase) UpdateCandidate(ctx context.Context, candidate *domain.Candidate) error {
	if err := u.candidateRepo.Update(ctx, candidate); err != nil {
		return fmt.Errorf("update candidate: %w", err)
	}

	return nil
}

func (u *candidateUseCase) UpdateCandidateProfile(ctx context.Context, profile *domain.CandidateProfile) error {
	if err := u.candidateRepo.UpdateProfile(ctx, profile); err != nil {
		return fmt.Errorf("update candidate profile: %w", err)
	}

	return nil
}

func (u *candidateUseCase) DeleteCandidate(ctx context.Context, id string) error {
	if err := u.candidateRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete candidate: %w", err)
	}

	return nil
}

func (u *candidateUseCase) MoveCandidate(ctx context.Context, p domain.MoveCandidateParams) error {
	if err := u.candidateRepo.MoveStage(ctx, p); err != nil {
		return fmt.Errorf("move candidate: %w", err)
	}

	return nil
}
