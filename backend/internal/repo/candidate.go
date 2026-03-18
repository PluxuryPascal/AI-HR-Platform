package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

type CandidateRepository interface {
	Create(ctx context.Context, jobID, fileKey string) (*domain.Candidate, error)
	GetByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error)
	GetByJobID(ctx context.Context, jobID string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error)
	UpdateByRecruiter(ctx context.Context, candidate *domain.Candidate) error
	UpdateFromAIParsing(ctx context.Context, result *domain.AIParsingResult) error
	UpdateParsingStatus(ctx context.Context, candidateID string, status domain.CandidateParsingStatus) error
	Delete(ctx context.Context, id string) error
	MoveStage(ctx context.Context, p domain.MoveCandidateParams) error
	GetStageHistory(ctx context.Context, candidateID string) ([]domain.StageHistoryEntry, error)

	// AI Support
	SaveCandidateScore(ctx context.Context, score *domain.CandidateScore, factors []domain.ScoreFactor) error
	GetScoreByCandidateID(ctx context.Context, candidateID string) (*domain.CandidateScore, []domain.ScoreFactor, error)
	GetScoresByCandidateIDs(ctx context.Context, candidateIDs []string) (map[string]*domain.CandidateScore, error)
	SaveResumeEmbedding(ctx context.Context, embedding *domain.ResumeEmbedding) error

}

type candidateRepo struct {
	dbClient *db.PostgresClient
}

func NewCandidateRepo(dbClient *db.PostgresClient) CandidateRepository {
	return &candidateRepo{dbClient: dbClient}
}

func (r *candidateRepo) UpdateParsingStatus(ctx context.Context, candidateID string, status domain.CandidateParsingStatus) error {
	const query = `
		UPDATE hiring.t_candidates
		SET 
			parsing_status = @parsing_status,
			updated_at = NOW()
		WHERE id = @id
	`

	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"id":             candidateID,
		"parsing_status": status,
	})
	if err != nil {
		return fmt.Errorf("update parsing status: %w", err)
	}

	return nil
}

func (r *candidateRepo) UpdateFromAIParsing(ctx context.Context, result *domain.AIParsingResult) error {
	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const query = `
		UPDATE hiring.t_candidates
		SET 
			first_name = @first_name,
			last_name = @last_name,
			email = @email,
			parsed_text = @parsed_text,
			location = @location,
			skills = @skills,
			parsing_status = @parsing_status,
			updated_at = NOW()
		WHERE id = @id
	`

	_, err = tx.Exec(ctx, query, pgx.NamedArgs{
		"id":             result.CandidateID,
		"first_name":     result.FirstName,
		"last_name":      result.LastName,
		"email":          result.Email,
		"parsed_text":    result.ParsedText,
		"location":       result.Location,
		"skills":         result.Skills,
		"parsing_status": result.ParseStatus,
	})
	if err != nil {
		return fmt.Errorf("update candidate: %w", err)
	}

	now := time.Now()
	const profileQuery = `
		INSERT INTO hiring.t_candidate_profiles (candidate_id, structured_data, missing_fields, ai_parsed_at, updated_at)
		VALUES (@candidate_id, @structured_data, @missing_fields, @ai_parsed_at, NOW())
		ON CONFLICT (candidate_id) DO UPDATE SET
			structured_data = EXCLUDED.structured_data,
			missing_fields = EXCLUDED.missing_fields,
			ai_parsed_at = EXCLUDED.ai_parsed_at,
			updated_at = NOW()
	`

	_, err = tx.Exec(ctx, profileQuery, pgx.NamedArgs{
		"candidate_id":    result.CandidateID,
		"structured_data": result.StructuredData,
		"missing_fields":  result.MissingFields,
		"ai_parsed_at":    now,
	})
	if err != nil {
		return fmt.Errorf("insert candidate profile: %w", err)
	}

	if result.InitialStageID != nil {
		const stageQuery = `
			INSERT INTO hiring.t_candidate_stages (candidate_id, stage_id, position, moved_at)
			VALUES (@candidate_id, @stage_id, 0, NOW())
			ON CONFLICT (candidate_id) DO NOTHING
		`

		_, err = tx.Exec(ctx, stageQuery, pgx.NamedArgs{
			"candidate_id": result.CandidateID,
			"stage_id":     *result.InitialStageID,
		})
		if err != nil {
			return fmt.Errorf("insert candidate stage: %w", err)
		}

		const historyQuery = `
			INSERT INTO hiring.t_candidate_stage_history (candidate_id, from_stage_id, to_stage_id, changed_by, changed_at)
			VALUES (@candidate_id, NULL, @stage_id, 'system', NOW())
		`

		_, err = tx.Exec(ctx, historyQuery, pgx.NamedArgs{
			"candidate_id": result.CandidateID,
			"stage_id":     *result.InitialStageID,
		})
		if err != nil {
			return fmt.Errorf("insert candidate stage history: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *candidateRepo) Create(ctx context.Context, jobID, fileKey string) (*domain.Candidate, error) {
	const query = `
		INSERT INTO hiring.t_candidates (job_id, resume_file_key, parsing_status)
		VALUES (@job_id, @resume_file_key, @parsing_status)
		RETURNING id, job_id, resume_file_key, parsing_status, created_at
	`
	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"job_id":          jobID,
		"resume_file_key": fileKey,
		"parsing_status":  domain.ParsingStatusPending,
	})
	if err != nil {
		return nil, fmt.Errorf("insert candidate draft: %w", err)
	}

	candidate, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Candidate])
	if err != nil {
		return nil, fmt.Errorf("collect candidate: %w", err)
	}

	return &candidate, nil
}

func (r *candidateRepo) GetByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error) {
	const query = `
		SELECT 
			c.id, c.job_id, c.first_name, c.last_name, c.email, c.resume_file_key, c.parsed_text, c.location, c.skills, c.parsing_status, c.created_at, c.updated_at,
			cp.structured_data, cp.missing_fields, cp.ai_parsed_at, cp.updated_at as profile_updated_at,
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
		&cand.ID, &cand.JobID, &cand.FirstName, &cand.LastName, &cand.Email, &cand.ResumeFileKey, &cand.ParsedText, &cand.Location, &cand.Skills, &cand.ParsingStatus, &cand.CreatedAt, &cand.UpdatedAt,
		&profile.StructuredData, &profile.MissingFields, &profile.AIParsedAt, &profile.UpdatedAt,
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
					'parsing_status',  l.parsing_status,
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

func (r *candidateRepo) UpdateByRecruiter(ctx context.Context, candidate *domain.Candidate) error {
	const query = `
		UPDATE hiring.t_candidates
		SET 
			first_name = @first_name, 
			last_name = @last_name, 
			email = @email,
		    location = @location, 
		    skills = @skills, 
		    updated_at = NOW()
		WHERE id = @id
		RETURNING updated_at
	`

	if err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id":         candidate.ID,
		"first_name": candidate.FirstName,
		"last_name":  candidate.LastName,
		"email":      candidate.Email,
		"location":   candidate.Location,
		"skills":     candidate.Skills,
	}).Scan(&candidate.UpdatedAt); err != nil {
		return fmt.Errorf("update candidate: %w", err)
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

func (r *candidateRepo) SaveCandidateScore(ctx context.Context, score *domain.CandidateScore, factors []domain.ScoreFactor) error {
	tx, err := r.dbClient.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const queryScore = `
		INSERT INTO ai_engine.t_candidate_scores (candidate_id, match_score, analyzed_at)
		VALUES (@candidate_id, @match_score, NOW())
		ON CONFLICT (candidate_id) DO UPDATE 
		SET match_score = EXCLUDED.match_score, analyzed_at = NOW()
	`
	_, err = tx.Exec(ctx, queryScore, pgx.NamedArgs{
		"candidate_id": score.CandidateID,
		"match_score":  score.MatchScore,
	})
	if err != nil {
		return fmt.Errorf("upsert candidate score: %w", err)
	}

	const queryDelFactors = `DELETE FROM ai_engine.t_score_factors WHERE candidate_id = @candidate_id`
	_, err = tx.Exec(ctx, queryDelFactors, pgx.NamedArgs{"candidate_id": score.CandidateID})
	if err != nil {
		return fmt.Errorf("delete old factors: %w", err)
	}

	if len(factors) > 0 {
		batch := &pgx.Batch{}
		const queryFact = `
			INSERT INTO ai_engine.t_score_factors (candidate_id, type, description, impact)
			VALUES ($1, $2, $3, $4)
		`
		for _, f := range factors {
			batch.Queue(queryFact, score.CandidateID, f.Type, f.Description, f.Impact)
		}
		br := tx.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return fmt.Errorf("insert factors batch: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *candidateRepo) GetScoreByCandidateID(ctx context.Context, candidateID string) (*domain.CandidateScore, []domain.ScoreFactor, error) {
	const queryScore = `
		SELECT candidate_id, match_score, analyzed_at
		FROM ai_engine.t_candidate_scores
		WHERE candidate_id = @id
	`
	var score domain.CandidateScore
	err := r.dbClient.Pool.QueryRow(ctx, queryScore, pgx.NamedArgs{"id": candidateID}).Scan(
		&score.CandidateID, &score.MatchScore, &score.AnalyzedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("query score: %w", err)
	}

	const queryFactors = `
		SELECT id, candidate_id, type, description, impact
		FROM ai_engine.t_score_factors
		WHERE candidate_id = @id
	`
	rows, err := r.dbClient.Pool.Query(ctx, queryFactors, pgx.NamedArgs{"id": candidateID})
	if err != nil {
		return nil, nil, fmt.Errorf("query factors: %w", err)
	}
	defer rows.Close()

	var factors []domain.ScoreFactor
	for rows.Next() {
		var f domain.ScoreFactor
		if err := rows.Scan(&f.ID, &f.CandidateID, &f.Type, &f.Description, &f.Impact); err != nil {
			return nil, nil, fmt.Errorf("scan factor: %w", err)
		}
		factors = append(factors, f)
	}

	return &score, factors, nil
}

func (r *candidateRepo) GetScoresByCandidateIDs(ctx context.Context, candidateIDs []string) (map[string]*domain.CandidateScore, error) {
	if len(candidateIDs) == 0 {
		return map[string]*domain.CandidateScore{}, nil
	}

	const queryScore = `
		SELECT candidate_id, match_score, analyzed_at
		FROM ai_engine.t_candidate_scores
		WHERE candidate_id = ANY(@ids)
	`
	rows, err := r.dbClient.Pool.Query(ctx, queryScore, pgx.NamedArgs{"ids": candidateIDs})
	if err != nil {
		return nil, fmt.Errorf("query scores: %w", err)
	}

	scores, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.CandidateScore])
	if err != nil {
		return nil, fmt.Errorf("scan scores: %w", err)
	}

	res := make(map[string]*domain.CandidateScore)
	for i := range scores {
		res[scores[i].CandidateID] = &scores[i]
	}

	return res, nil
}

func (r *candidateRepo) SaveResumeEmbedding(ctx context.Context, embedding *domain.ResumeEmbedding) error {
	const query = `
		INSERT INTO ai_engine.t_resume_embeddings (team_id, candidate_id, chunk_text, embedding)
		VALUES (@team_id, @candidate_id, @chunk_text, @embedding)
	`
	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"team_id":      embedding.TeamID,
		"candidate_id": embedding.CandidateID,
		"chunk_text":   embedding.ChunkText,
		"embedding":    pgvector.NewVector(embedding.Embedding),
	})
	if err != nil {
		return fmt.Errorf("insert embedding: %w", err)
	}
	return nil
}

func (r *candidateRepo) GetStageHistory(ctx context.Context, candidateID string) ([]domain.StageHistoryEntry, error) {
	const query = `
		SELECT 
			h.id,
			h.candidate_id,
			h.from_stage_id,
			fs.title AS from_stage_title,
			h.to_stage_id,
			ts.title AS to_stage_title,
			h.changed_by,
			h.changed_at
		FROM hiring.t_candidate_stage_history h
		LEFT JOIN hiring.t_pipeline_stages fs ON h.from_stage_id = fs.id
		JOIN hiring.t_pipeline_stages ts ON h.to_stage_id = ts.id
		WHERE h.candidate_id = @candidate_id
		ORDER BY h.changed_at DESC
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"candidate_id": candidateID})
	if err != nil {
		return nil, fmt.Errorf("query stage history: %w", err)
	}
	defer rows.Close()

	var entries []domain.StageHistoryEntry
	for rows.Next() {
		var e domain.StageHistoryEntry
		if err := rows.Scan(
			&e.ID, &e.CandidateID,
			&e.FromStageID, &e.FromStageTitle,
			&e.ToStageID, &e.ToStageTitle,
			&e.ChangedBy, &e.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("scan history entry: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, nil
}

