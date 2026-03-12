package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DepartmentRepository interface {
	Create(ctx context.Context, params *domain.CreateDepartmentParams) (*domain.Department, error)
	GetByID(ctx context.Context, id string) (*domain.Department, error)
	GetByTeamID(ctx context.Context, teamID string) ([]domain.Department, error)
	Delete(ctx context.Context, id string) error
}

type departmentRepo struct {
	dbClient *db.PostgresClient
}

func NewDepartmentRepo(dbClient *db.PostgresClient) DepartmentRepository {
	return &departmentRepo{dbClient: dbClient}
}

func (r *departmentRepo) Create(ctx context.Context, params *domain.CreateDepartmentParams) (*domain.Department, error) {
	const query = `
		INSERT INTO hiring.t_departments (team_id, name)
		VALUES (@team_id, @name)
		RETURNING id, team_id, name
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"team_id": params.TeamID,
		"name":    params.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("insert department: %w", err)
	}

	dept, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Department])
	if err != nil {
		return nil, fmt.Errorf("scan department: %w", err)
	}

	return &dept, nil
}

func (r *departmentRepo) GetByID(ctx context.Context, id string) (*domain.Department, error) {
	const query = `
		SELECT id, team_id, name
		FROM hiring.t_departments
		WHERE id = @id
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("query department by id: %w", err)
	}

	dept, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Department])
	if err != nil {
		return nil, fmt.Errorf("scan department: %w", err)
	}

	return &dept, nil
}

func (r *departmentRepo) GetByTeamID(ctx context.Context, teamID string) ([]domain.Department, error) {
	const query = `
		SELECT id, team_id, name
		FROM hiring.t_departments
		WHERE team_id = @team_id
		ORDER BY name ASC
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"team_id": teamID})
	if err != nil {
		return nil, fmt.Errorf("query departments by team id: %w", err)
	}

	depts, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Department])
	if err != nil {
		return nil, fmt.Errorf("scan departments: %w", err)
	}

	return depts, nil
}

func (r *departmentRepo) Delete(ctx context.Context, id string) error {
	const query = `
		DELETE FROM hiring.t_departments
		WHERE id = @id
	`

	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("delete department: %w", err)
	}

	return nil
}
