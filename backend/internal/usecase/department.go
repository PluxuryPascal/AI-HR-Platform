package usecase

import (
	"backend/internal/domain"
	"backend/internal/repo"
	"context"
	"fmt"
)

type DepartmentUseCase interface {
	CreateDepartment(ctx context.Context, params *domain.CreateDepartmentParams) (*domain.Department, error)
	GetDepartmentsByTeam(ctx context.Context, teamID string) ([]domain.Department, error)
	DeleteDepartment(ctx context.Context, id string) error
}

type departmentUseCase struct {
	deptRepo repo.DepartmentRepository
}

func NewDepartmentUseCase(deptRepo repo.DepartmentRepository) DepartmentUseCase {
	return &departmentUseCase{deptRepo: deptRepo}
}

func (u *departmentUseCase) CreateDepartment(ctx context.Context, params *domain.CreateDepartmentParams) (*domain.Department, error) {
	dept, err := u.deptRepo.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create department: %w", err)
	}
	return dept, nil
}

func (u *departmentUseCase) GetDepartmentsByTeam(ctx context.Context, teamID string) ([]domain.Department, error) {
	depts, err := u.deptRepo.GetByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("get departments by team: %w", err)
	}
	return depts, nil
}

func (u *departmentUseCase) DeleteDepartment(ctx context.Context, id string) error {
	if err := u.deptRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete department: %w", err)
	}
	return nil
}
