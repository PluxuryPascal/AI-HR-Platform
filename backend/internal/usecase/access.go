package usecase

import (
	"backend/internal/repo"
	"context"
	"fmt"
)

type AccessUseCase interface {
	TransferAccess(ctx context.Context, userID string, jobIDs []string) error
}

type accessUseCase struct {
	repo repo.AccessRepository
}

func NewAccessUseCase(repo repo.AccessRepository) AccessUseCase {
	return &accessUseCase{repo: repo}
}

func (u *accessUseCase) TransferAccess(ctx context.Context, userID string, jobIDs []string) error {
	if err := u.repo.TransferAccess(ctx, userID, jobIDs); err != nil {
		return fmt.Errorf("transfer access error: %w", err)
	}

	return nil
}
