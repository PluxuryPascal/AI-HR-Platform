package usecase

import (
	"backend/internal/domain"
	"backend/internal/repo"
	"context"
	"fmt"
)

type CandidateUseCase interface {
	CreateCandidate(ctx context.Context, candidate *domain.Candidate, profile *domain.CandidateProfile, initialStageID string) error
	GetCandidateByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error)
	GetCandidatesByJob(ctx context.Context, jobID string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error)
	UpdateCandidate(ctx context.Context, candidate *domain.Candidate) error
	UpdateCandidateProfile(ctx context.Context, profile *domain.CandidateProfile) error
	DeleteCandidate(ctx context.Context, id string) error
}

type candidateUseCase struct {
	candidateRepo repo.CandidateRepository
}

func NewCandidateUseCase(candidateRepo repo.CandidateRepository) CandidateUseCase {
	return &candidateUseCase{candidateRepo: candidateRepo}
}

func (u *candidateUseCase) CreateCandidate(ctx context.Context, candidate *domain.Candidate, profile *domain.CandidateProfile, initialStageID string) error {
	if err := u.candidateRepo.Create(ctx, candidate, profile, initialStageID); err != nil {
		return fmt.Errorf("create candidate: %w", err)
	}

	return nil
}

func (u *candidateUseCase) GetCandidateByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error) {
	cand, profile, stageID, err := u.candidateRepo.GetByID(ctx, id)
	if err != nil {
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
