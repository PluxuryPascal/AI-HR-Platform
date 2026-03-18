package usecase

import (
	"backend/internal/audit"
	"backend/internal/domain"
	"backend/internal/repo"
	"context"
	"fmt"
)

type AccessUseCase interface {
	TransferAccess(ctx context.Context, userID string, jobIDs []string) error
	GrantAccess(ctx context.Context, userID, jobID, actorID, teamID string) error
	RevokeAccess(ctx context.Context, userID, jobID, actorID, teamID string) error
	GetAccessByJobID(ctx context.Context, jobID string) ([]domain.JobAccess, error)
}

type accessUseCase struct {
	repo    repo.AccessRepository
	auditor *audit.Logger
}

func NewAccessUseCase(repo repo.AccessRepository, auditor *audit.Logger) AccessUseCase {
	return &accessUseCase{
		repo:    repo,
		auditor: auditor,
	}
}

func (u *accessUseCase) TransferAccess(ctx context.Context, userID string, jobIDs []string) error {
	if err := u.repo.TransferAccess(ctx, userID, jobIDs); err != nil {
		return fmt.Errorf("transfer access error: %w", err)
	}
	return nil
}

func (u *accessUseCase) GrantAccess(ctx context.Context, userID, jobID, actorID, teamID string) error {
	if err := u.repo.GrantAccess(ctx, userID, jobID); err != nil {
		return fmt.Errorf("grant access: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    teamID,
		ActorType: audit.ActorUser,
		ActorID:   &actorID,
		Action:    audit.HiringAccessGranted,
		TargetID:  &jobID,
		Payload:   map[string]string{"granted_to": userID},
	})

	return nil
}

func (u *accessUseCase) RevokeAccess(ctx context.Context, userID, jobID, actorID, teamID string) error {
	if err := u.repo.RevokeAccess(ctx, userID, jobID); err != nil {
		return fmt.Errorf("revoke access: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    teamID,
		ActorType: audit.ActorUser,
		ActorID:   &actorID,
		Action:    audit.HiringAccessRevoked,
		TargetID:  &jobID,
		Payload:   map[string]string{"revoked_from": userID},
	})

	return nil
}

func (u *accessUseCase) GetAccessByJobID(ctx context.Context, jobID string) ([]domain.JobAccess, error) {
	list, err := u.repo.GetAccessByJobID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get access by job: %w", err)
	}
	return list, nil
}

