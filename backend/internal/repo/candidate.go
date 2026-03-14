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

type CandidateRepository interface {
	Create(ctx context.Context, candidate *domain.Candidate, profile *domain.CandidateProfile, initialStageID string) error
	GetByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error)
	GetByJobID(ctx context.Context, jobID string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error)
	Update(ctx context.Context, candidate *domain.Candidate) error
	UpdateProfile(ctx context.Context, profile *domain.CandidateProfile) error
	Delete(ctx context.Context, id string) error
	MoveStage(ctx context.Context, p domain.MoveCandidateParams) error
}

type candidateRepo struct {
	dbClient *db.PostgresClient
}

func NewCandidateRepo(dbClient *db.PostgresClient) CandidateRepository {
	return &candidateRepo{dbClient: dbClient}
}

func (r *candidateRepo) Create(ctx context.Context, candidate *domain.Candidate, profile *domain.CandidateProfile, initialStageID string) error {
	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	// 1. Insert Candidate
	const queryCand = `
		INSERT INTO hiring.t_candidates (job_id, first_name, last_name, email, resume_file_key, parsed_text, location, skills)
		VALUES (@job_id, @first_name, @last_name, @email, @resume_file_key, @parsed_text, @location, @skills)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, queryCand, pgx.NamedArgs{
		"job_id":          candidate.JobID,
		"first_name":      candidate.FirstName,
		"last_name":       candidate.LastName,
		"email":           candidate.Email,
		"resume_file_key": candidate.ResumeFileKey,
		"parsed_text":     candidate.ParsedText,
		"location":        candidate.Location,
		"skills":          candidate.Skills,
	}).Scan(&candidate.ID, &candidate.CreatedAt, &candidate.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert candidate: %w", err)
	}

	// 2. Insert Profile
	const queryProfile = `
		INSERT INTO hiring.t_candidate_profiles (candidate_id, structured_data, ai_parsed_at)
		VALUES (@candidate_id, @structured_data, @ai_parsed_at)
	`
	_, err = tx.Exec(ctx, queryProfile, pgx.NamedArgs{
		"candidate_id":    candidate.ID,
		"structured_data": profile.StructuredData,
		"ai_parsed_at":    profile.AIParsedAt,
	})
	if err != nil {
		return fmt.Errorf("insert candidate profile: %w", err)
	}

	// 3. Set Initial Stage
	const queryStage = `
		INSERT INTO hiring.t_candidate_stages (candidate_id, stage_id, position)
		VALUES (@candidate_id, @stage_id, 0)
	`
	_, err = tx.Exec(ctx, queryStage, pgx.NamedArgs{
		"candidate_id": candidate.ID,
		"stage_id":     initialStageID,
	})
	if err != nil {
		return fmt.Errorf("insert candidate stage: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *candidateRepo) GetByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error) {
	const query = `
		SELECT 
			c.id, c.job_id, c.first_name, c.last_name, c.email, c.resume_file_key, c.parsed_text, c.location, c.skills, c.created_at, c.updated_at,
			cp.structured_data, cp.ai_parsed_at, cp.updated_at as profile_updated_at,
			cs.stage_id
		FROM hiring.t_candidates c
		LEFT JOIN hiring.t_candidate_profiles cp ON c.id = cp.candidate_id
		LEFT JOIN hiring.t_candidate_stages cs ON c.id = cs.candidate_id
		WHERE c.id = @id
	`

	var cand domain.Candidate
	var profile domain.CandidateProfile
	var stageID *string

	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{"id": id}).Scan(
		&cand.ID, &cand.JobID, &cand.FirstName, &cand.LastName, &cand.Email, &cand.ResumeFileKey, &cand.ParsedText, &cand.Location, &cand.Skills, &cand.CreatedAt, &cand.UpdatedAt,
		&profile.StructuredData, &profile.AIParsedAt, &profile.UpdatedAt,
		&stageID,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("query candidate: %w", err)
	}

	profile.CandidateID = cand.ID
	return &cand, &profile, stageID, nil
}

func (r *candidateRepo) GetByJobID(ctx context.Context, jobID string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error) {
	sortID := "c.created_at"
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
			"first_name": "c.first_name",
			"last_name":  "c.last_name",
			"created_at": "c.created_at",
		}

		if dbField, ok := allowedSortFields[sortID]; ok {
			sortID = dbField
		} else {
			sortID = "c.created_at"
		}
	}

	rawQuery := fmt.Sprintf(`
		WITH
			all_candidates AS (
				SELECT c.*, cs.stage_id
				FROM hiring.t_candidates c
				LEFT JOIN hiring.t_candidate_stages cs ON c.id = cs.candidate_id
				WHERE c.job_id = @job_id
			),
			
			filtered AS (
				SELECT ac.*
				FROM all_candidates ac
				WHERE
					CASE 
						WHEN @ifFirstName THEN 
							lower(ac.first_name) ILIKE @first_name
						ELSE 
							TRUE 
					END
					AND 
					CASE 
						WHEN @ifLastName THEN 
							lower(ac.last_name) ILIKE @last_name
						ELSE 
							TRUE 
					END
					AND
					CASE 
						WHEN @ifEmail THEN 
							lower(ac.email) ILIKE @email
						ELSE 
							TRUE 
					END
					AND
					CASE 
						WHEN @ifStage THEN 
							ac.stage_id = @stage_id
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
					'id',              l.id,
					'job_id',          l.job_id,
					'first_name',      l.first_name,
					'last_name',       l.last_name,
					'email',           l.email,
					'resume_file_key', l.resume_file_key,
					'location',        l.location,
					'skills',          l.skills,
					'created_at',      l.created_at,
					'updated_at',      l.updated_at
				)) AS build
				FROM limited l
			)
			
		SELECT
			(SELECT count(1) FROM filtered) AS total,
			(
				SELECT COALESCE(b.build, '[]'::JSONB)
				FROM build b
			) AS candidates;
	`, sortID, sortDesc)

	args := pgx.NamedArgs{
		"job_id": jobID,
		"offset": offset,
		"limit":  limit,
	}

	if filter.FirstName != nil {
		args["ifFirstName"] = true
		args["first_name"] = "%" + strings.ToLower(*filter.FirstName) + "%"
	} else {
		args["ifFirstName"] = false
		args["first_name"] = ""
	}

	if filter.LastName != nil {
		args["ifLastName"] = true
		args["last_name"] = "%" + strings.ToLower(*filter.LastName) + "%"
	} else {
		args["ifLastName"] = false
		args["last_name"] = ""
	}

	if filter.Email != nil {
		args["ifEmail"] = true
		args["email"] = "%" + strings.ToLower(*filter.Email) + "%"
	} else {
		args["ifEmail"] = false
		args["email"] = ""
	}

	if filter.CurrentStageID != nil {
		args["ifStage"] = true
		args["stage_id"] = *filter.CurrentStageID
	} else {
		args["ifStage"] = false
		args["stage_id"] = ""
	}

	var dateFilterClause string
	if filter.DateFilter != nil {
		dateType := strings.ToLower(filter.DateFilter.Type)
		switch dateType {
		case "before":
			if filter.DateFilter.DateTo != nil {
				dateFilterClause = "AND ac.created_at <= @date_to"
				args["ifDate"] = true
				args["date_to"] = filter.DateFilter.DateTo
			} else {
				args["ifDate"] = false
			}
		case "after":
			if filter.DateFilter.DateFrom != nil {
				dateFilterClause = "AND ac.created_at >= @date_from"
				args["ifDate"] = true
				args["date_from"] = filter.DateFilter.DateFrom
			} else {
				args["ifDate"] = false
			}
		case "between":
			if filter.DateFilter.DateFrom != nil && filter.DateFilter.DateTo != nil {
				dateFilterClause = "AND ac.created_at BETWEEN @date_from AND @date_to"
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
		return nil, fmt.Errorf("query candidates: %w", err)
	}
	defer rows.Close()

	var res domain.CandidatesDTO
	if rows.Next() {
		err = rows.Scan(&res.Total, &res.Candidates)
		if err != nil {
			return nil, fmt.Errorf("scan candidates dto: %w", err)
		}
	}

	return &res, nil
}

func (r *candidateRepo) Update(ctx context.Context, candidate *domain.Candidate) error {
	const query = `
		UPDATE hiring.t_candidates
		SET first_name = @first_name, last_name = @last_name, email = @email,
		    resume_file_key = @resume_file_key, parsed_text = @parsed_text, 
		    location = @location, skills = @skills, updated_at = NOW()
		WHERE id = @id
		RETURNING updated_at
	`
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id":              candidate.ID,
		"first_name":      candidate.FirstName,
		"last_name":       candidate.LastName,
		"email":           candidate.Email,
		"resume_file_key": candidate.ResumeFileKey,
		"parsed_text":     candidate.ParsedText,
		"location":        candidate.Location,
		"skills":          candidate.Skills,
	}).Scan(&candidate.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update candidate: %w", err)
	}
	return nil
}

func (r *candidateRepo) UpdateProfile(ctx context.Context, profile *domain.CandidateProfile) error {
	const query = `
		UPDATE hiring.t_candidate_profiles
		SET structured_data = @structured_data, ai_parsed_at = @ai_parsed_at, updated_at = NOW()
		WHERE candidate_id = @candidate_id
		RETURNING updated_at
	`
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"candidate_id":    profile.CandidateID,
		"structured_data": profile.StructuredData,
		"ai_parsed_at":    profile.AIParsedAt,
	}).Scan(&profile.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update candidate profile: %w", err)
	}
	return nil
}

func (r *candidateRepo) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM hiring.t_candidates WHERE id = @id`
	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("delete candidate: %w", err)
	}
	return nil
}

func (r *candidateRepo) MoveStage(ctx context.Context, p domain.MoveCandidateParams) error {
	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Get current stage
	const queryGetCurr = `SELECT stage_id FROM hiring.t_candidate_stages WHERE candidate_id = @candidate_id`
	var fromStageID *string
	err = tx.QueryRow(ctx, queryGetCurr, pgx.NamedArgs{"candidate_id": p.CandidateID}).Scan(&fromStageID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("get current stage: %w", err)
	}

	// 2. Update current stage
	const queryUpdateStage = `
		INSERT INTO hiring.t_candidate_stages (candidate_id, stage_id, position, moved_at)
		VALUES (@candidate_id, @stage_id, @position, NOW())
		ON CONFLICT (candidate_id) DO UPDATE 
		SET stage_id = EXCLUDED.stage_id, position = EXCLUDED.position, moved_at = NOW()
	`
	_, err = tx.Exec(ctx, queryUpdateStage, pgx.NamedArgs{
		"candidate_id": p.CandidateID,
		"stage_id":     p.ToStageID,
		"position":     p.NewPosition,
	})
	if err != nil {
		return fmt.Errorf("update candidate stage: %w", err)
	}

	// 3. Record history
	const queryHistory = `
		INSERT INTO hiring.t_candidate_stage_history (candidate_id, from_stage_id, to_stage_id, changed_by)
		VALUES (@candidate_id, @from_stage_id, @to_stage_id, @changed_by)
	`
	_, err = tx.Exec(ctx, queryHistory, pgx.NamedArgs{
		"candidate_id":  p.CandidateID,
		"from_stage_id": fromStageID,
		"to_stage_id":   p.ToStageID,
		"changed_by":    p.ChangedBy,
	})
	if err != nil {
		return fmt.Errorf("insert stage history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
