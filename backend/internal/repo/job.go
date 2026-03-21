package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

type JobRepository interface {
	Create(ctx context.Context, job *domain.Job) error
	GetByID(ctx context.Context, id string) (*domain.Job, error)
	GetByTeamID(ctx context.Context, teamID string, offset, limit int, filter domain.JobFilter) (*domain.JobsDTO, error)
	Update(ctx context.Context, job *domain.Job) error
	Delete(ctx context.Context, id string) error
	GetJobRequirements(ctx context.Context, jobID string) ([]byte, error)
}

type jobRepo struct {
	dbClient *db.PostgresClient
}

func NewJobRepo(dbClient *db.PostgresClient) JobRepository {
	return &jobRepo{dbClient: dbClient}
}

func (r *jobRepo) Create(ctx context.Context, job *domain.Job) error {
	const query = `
		INSERT INTO hiring.t_jobs (
			team_id, title, department_id, work_format, description, 
			extracted_requirements, status, salary_min, salary_max, currency, created_by
		) VALUES (
			@team_id, @title, @department_id, @work_format, @description, 
			@extracted_requirements, @status, @salary_min, @salary_max, @currency, @created_by
		) RETURNING id, created_at, updated_at
	`

	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"team_id":                job.TeamID,
		"title":                  job.Title,
		"department_id":          job.DepartmentID,
		"work_format":            job.WorkFormat,
		"description":            job.Description,
		"extracted_requirements": job.ExtractedRequirements,
		"status":                 job.Status,
		"salary_min":             job.SalaryMin,
		"salary_max":             job.SalaryMax,
		"currency":               job.Currency,
		"created_by":             job.CreatedBy,
	}).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil
}

func (r *jobRepo) GetByID(ctx context.Context, id string) (*domain.Job, error) {
	const query = `
		SELECT id, team_id, title, department_id, work_format, description, 
		       extracted_requirements, status, salary_min, salary_max, currency, 
		       created_by, created_at, updated_at
		FROM hiring.t_jobs
		WHERE id = @id
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("query job by id: %w", err)
	}

	job, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Job])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scan job: %w", err)
	}

	return &job, nil
}

func (r *jobRepo) GetByTeamID(ctx context.Context, teamID string, offset, limit int, filter domain.JobFilter) (*domain.JobsDTO, error) {
	sortID := "f.created_at"
	sortDesc := "DESC"

	if filter.Sort != nil {
		if filter.Sort.SortID != nil {
			sortID = *filter.Sort.SortID
		}

		if filter.Sort.SortDesc != nil {
			if *filter.Sort.SortDesc {
				sortDesc = "DESC"
			} else {
				sortDesc = "ASC"
			}
		}

		allowedSortFields := map[string]string{
			"title":           "f.title",
			"created_at":      "f.created_at",
			"status":          "f.status",
			"department_name": "f.department_name",
		}

		if dbField, ok := allowedSortFields[sortID]; ok {
			sortID = dbField
		} else {
			sortID = "f.created_at"
		}
	}

	rawQuery := fmt.Sprintf(`
		WITH
			all_jobs AS (
				SELECT j.*, d.name as department_name
				FROM hiring.t_jobs j
				LEFT JOIN hiring.t_departments d ON j.department_id = d.id
				LEFT JOIN hiring.t_job_access ja ON j.id = ja.job_id
				WHERE j.team_id = @team_id
				  AND (
				    @ifAllowedUser = FALSE 
				    OR ja.user_id = @allowed_user_id
				  )
			),
			
			filtered AS (
				SELECT aj.*
				FROM all_jobs aj
				WHERE
					CASE 
						WHEN @ifTitle THEN 
							lower(aj.title) ILIKE @title
						ELSE 
							TRUE 
					END
					
					AND 
					
					CASE 
						WHEN @ifStatus THEN 
							aj.status = @status::job_status
						ELSE 
							TRUE 
					END
					
					AND

					CASE 
						WHEN @ifWorkFormat THEN 
							aj.work_format = @work_format::work_format_type
						ELSE 
							TRUE 
					END

					AND

					CASE 
						WHEN @ifDepartmentName THEN 
							lower(aj.department_name) ILIKE @department_name
						ELSE 
							TRUE 
					END

					AND
					
					CASE 
						WHEN @ifDate THEN 
							TRUE
							{{DATE_FILTER}}
						ELSE 
							TRUE 
					END			
			),
			
			limited AS (
				SELECT f.*
				FROM filtered f
				ORDER BY %s %s
				LIMIT @limit
				OFFSET @offset
			),
			
			build AS (
				SELECT jsonb_agg(jsonb_build_object(
					'id',                    l.id,
					'team_id',               l.team_id,
					'title',                 l.title,
					'department_id',         l.department_id,
					'department_name',       l.department_name,
					'work_format',           l.work_format,
					'description',           l.description,
					'extracted_requirements',l.extracted_requirements,
					'status',                l.status,
					'salary_min',            l.salary_min,
					'salary_max',            l.salary_max,
					'currency',              l.currency,
					'created_by',            l.created_by,
					'created_at',            l.created_at,
					'updated_at',            l.updated_at
				)) AS build
				FROM limited l
			)
			
		SELECT
			(SELECT count(1) FROM filtered) AS total,
			(
				SELECT COALESCE(b.build, '[]'::JSONB)
				FROM build b
			) AS jobs;
	`, sortID, sortDesc)

	args := pgx.NamedArgs{
		"team_id": teamID,
		"offset":  offset,
		"limit":   limit,
	}

	if filter.AllowedUserID != nil {
		args["ifAllowedUser"] = true
		args["allowed_user_id"] = *filter.AllowedUserID
	} else {
		args["ifAllowedUser"] = false
		args["allowed_user_id"] = nil
	}

	if filter.Title != nil {

		args["ifTitle"] = true
		args["title"] = "%" + strings.ToLower(*filter.Title) + "%"
	} else {
		args["ifTitle"] = false
		args["title"] = ""
	}

	if filter.Status != nil {
		args["ifStatus"] = true
		args["status"] = string(*filter.Status)
	} else {
		args["ifStatus"] = false
		args["status"] = nil
	}

	if filter.WorkFormat != nil {
		args["ifWorkFormat"] = true
		args["work_format"] = string(*filter.WorkFormat)
	} else {
		args["ifWorkFormat"] = false
		args["work_format"] = nil
	}

	if filter.DepartmentName != nil {
		args["ifDepartmentName"] = true
		args["department_name"] = "%" + strings.ToLower(*filter.DepartmentName) + "%"
	} else {
		args["ifDepartmentName"] = false
		args["department_name"] = ""
	}

	var dateFilterClause string
	if filter.DateFilter != nil {
		dateType := strings.ToLower(filter.DateFilter.Type)

		switch dateType {
		case "before":
			if filter.DateFilter.DateTo != nil {
				dateFilterClause = "AND aj.created_at <= @date_to"
				args["ifDate"] = true
				args["date_to"] = filter.DateFilter.DateTo
			} else {
				args["ifDate"] = false
			}
		case "after":
			if filter.DateFilter.DateFrom != nil {
				dateFilterClause = "AND aj.created_at >= @date_from"
				args["ifDate"] = true
				args["date_from"] = filter.DateFilter.DateFrom
			} else {
				args["ifDate"] = false
			}
		case "between":
			if filter.DateFilter.DateFrom != nil && filter.DateFilter.DateTo != nil {
				dateFilterClause = "AND aj.created_at BETWEEN @date_from AND @date_to"
				args["ifDate"] = true
				args["date_from"] = filter.DateFilter.DateFrom
				args["date_to"] = filter.DateFilter.DateTo
			} else {
				args["ifDate"] = false
			}
		}
	} else {
		args["ifDate"] = false
	}

	finalQuery := strings.Replace(rawQuery, "{{DATE_FILTER}}", dateFilterClause, 1)

	rows, err := r.dbClient.Pool.Query(ctx, finalQuery, args)
	if err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}

	jobsRes, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.JobsDTO])
	if err != nil {
		return nil, fmt.Errorf("collect rows error: %w", err)
	}

	return &jobsRes, nil
}

func (r *jobRepo) Update(ctx context.Context, job *domain.Job) error {
	const query = `
		UPDATE hiring.t_jobs
		SET title = @title, department_id = @department_id, work_format = @work_format, 
		    description = @description, extracted_requirements = @extracted_requirements, 
		    status = @status, salary_min = @salary_min, salary_max = @salary_max, 
		    currency = @currency, updated_at = NOW()
		WHERE id = @id
		RETURNING updated_at
	`

	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id":                     job.ID,
		"title":                  job.Title,
		"department_id":          job.DepartmentID,
		"work_format":            job.WorkFormat,
		"description":            job.Description,
		"extracted_requirements": job.ExtractedRequirements,
		"status":                 job.Status,
		"salary_min":             job.SalaryMin,
		"salary_max":             job.SalaryMax,
		"currency":               job.Currency,
	}).Scan(&job.UpdatedAt)

	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}

	return nil
}

func (r *jobRepo) Delete(ctx context.Context, id string) error {
	const query = `
		DELETE FROM hiring.t_jobs
		WHERE id = @id
	`

	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("delete job: %w", err)
	}

	return nil
}

func (r *jobRepo) GetJobRequirements(ctx context.Context, jobID string) ([]byte, error) {
	const query = `
		SELECT extracted_requirements
		FROM hiring.t_jobs
		WHERE id = @id AND extracted_requirements IS NOT NULL
	`

	var requirements []byte
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{"id": jobID}).Scan(&requirements)
	if err != nil {
		return nil, err
	}

	return requirements, nil
}
