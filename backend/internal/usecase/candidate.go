package usecase

import (
	"backend/internal/audit"
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/pkg/mq"
	"backend/pkg/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrCandidateNotFound       = errors.New("candidate not found")
	ErrStorageUnavailable      = errors.New("storage unavailable")
	ErrCandidateNotNeedsReview = errors.New("candidate is not in needs review")
	ErrInvalidStageTransition  = errors.New("invalid stage transition: only +1 position or terminal stage allowed")
)

type CandidateUseCase interface {
	CreateCandidate(ctx context.Context, params domain.CreateCandidateParams) (*domain.Candidate, error)
	GetCandidateByID(ctx context.Context, id, userID, role string) (*domain.Candidate, *domain.CandidateProfile, *string, *domain.CandidateScore, []domain.ScoreFactor, error)
	GetCandidatesByJob(ctx context.Context, jobID, userID, role string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error)
	DeleteCandidate(ctx context.Context, id, actorID, teamID, role string) error
	MoveCandidate(ctx context.Context, p domain.MoveCandidateParams, role string) error
	BulkMoveCandidates(ctx context.Context, p domain.BulkMoveCandidateParams, role string) error
	FinalizeAIParsing(ctx context.Context, result domain.AIParsingResult) error
	UpdateByRecruiter(ctx context.Context, candidate *domain.Candidate, userID, role string) error
	ConfirmManualReview(ctx context.Context, candidate *domain.Candidate, userID, role string) error
	SaveInterviewGuide(ctx context.Context, candidateID string, guide []byte) error
	GetStageHistory(ctx context.Context, candidateID, userID, role string) ([]domain.StageHistoryEntry, error)
}

type candidateUseCase struct {
	log           *zap.Logger
	candidateRepo repo.CandidateRepository
	pipelineRepo  repo.PipelineRepository
	jobRepo       repo.JobRepository
	accessRepo    repo.AccessRepository
	storage       storage.FileStorage
	publisher     *mq.MQPublisher
	auditor       *audit.Logger
}

func (u *candidateUseCase) ConfirmManualReview(ctx context.Context, candidate *domain.Candidate, userID, role string) error {
	current, _, _, _, _, err := u.GetCandidateByID(ctx, candidate.ID, userID, role)
	if err != nil {
		return fmt.Errorf("get candidate by id: %w", err)
	}

	if current.ParsingStatus != domain.ParsingStatusNeedsReview {
		return ErrCandidateNotNeedsReview
	}

	candidate.ParsingStatus = domain.ParsingStatusCompleted

	if err := u.candidateRepo.UpdateByRecruiter(ctx, candidate); err != nil {
		return fmt.Errorf("update candidate by recruiter: %w", err)
	}

	return nil
}

func (u *candidateUseCase) SaveInterviewGuide(ctx context.Context, candidateID string, guide []byte) error {
	return u.candidateRepo.SaveInterviewGuide(ctx, candidateID, guide)
}

func (u *candidateUseCase) FinalizeAIParsing(ctx context.Context, result domain.AIParsingResult) error {
	if result.ParseStatus != domain.ParsingStatusFailed {
		stages, err := u.pipelineRepo.GetStagesByJobID(ctx, result.JobID)
		if err != nil {
			return fmt.Errorf("get stages: %w", err)
		}

		if len(stages) > 0 {
			result.InitialStageID = &stages[0].ID
		}

	}

	if err := u.candidateRepo.UpdateFromAIParsing(ctx, &result); err != nil {
		return fmt.Errorf("update candidate from ai parsing: %w", err)
	}

	if result.ParseStatus == domain.ParsingStatusNeedsReview {
		u.log.Info("candidate requires manual review",
			zap.String("candidate_id", result.CandidateID),
			zap.Strings("missing_fields", result.MissingFields),
		)
	}

	return nil
}

func (u *candidateUseCase) UpdateByRecruiter(ctx context.Context, candidate *domain.Candidate, userID, role string) error {
	if _, _, _, _, _, err := u.GetCandidateByID(ctx, candidate.ID, userID, role); err != nil {
		return err
	}

	if err := u.candidateRepo.UpdateByRecruiter(ctx, candidate); err != nil {
		return fmt.Errorf("update candidate by recruiter: %w", err)
	}

	return nil
}

func (u *candidateUseCase) CreateCandidate(ctx context.Context, params domain.CreateCandidateParams) (*domain.Candidate, error) {
	// Get job to extract TeamID for the event
	job, err := u.jobRepo.GetByID(ctx, params.JobID)
	if err != nil {
		return nil, fmt.Errorf("get job for team_id: %w", err)
	}

	fileKey, err := u.storage.UploadFile(ctx, params.File, params.Filename)
	if err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	candidate, err := u.candidateRepo.Create(ctx, params.JobID, fileKey)
	if err != nil {
		return nil, fmt.Errorf("create candidate: %w", err)
	}

	payload, err := json.Marshal(domain.CandidateCreatedEvent{
		CandidateID:   candidate.ID,
		JobID:         candidate.JobID,
		ResumeFileKey: *candidate.ResumeFileKey,
		TeamID:        job.TeamID,
		Locale:        params.Locale,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal candidate created event: %w", err)
	}

	if err := u.publisher.Publish(ctx, payload,
		mq.WithExchange("hiring.events"),
		mq.WithRoutingKey("candidate.created"),
	); err != nil {
		return nil, fmt.Errorf("publish candidate created event: %w", err)
	}

	candidate.ParsingStatus = domain.ParsingStatusProcessing
	if err := u.candidateRepo.UpdateParsingStatus(ctx, candidate.ID, candidate.ParsingStatus); err != nil {
		u.log.Warn("update parsing status", zap.String("candidate_id", candidate.ID), zap.Error(err))
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    job.TeamID,
		ActorType: audit.ActorUser,
		ActorID:   &params.ActorID,
		Action:    audit.HiringCandidateAdded,
		TargetID:  &candidate.ID,
		Payload:   map[string]string{"job_id": params.JobID},
	})

	return candidate, nil
}

func (u *candidateUseCase) GetCandidateByID(ctx context.Context, id, userID, role string) (*domain.Candidate, *domain.CandidateProfile, *string, *domain.CandidateScore, []domain.ScoreFactor, error) {
	cand, profile, stageID, err := u.candidateRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, nil, nil, nil, ErrCandidateNotFound
		}

		return nil, nil, nil, nil, nil, fmt.Errorf("get candidate by id: %w", err)
	}

	if role != "owner" && role != "admin" {
		ok, err := u.accessRepo.HasAccess(ctx, userID, cand.JobID)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("check job access: %w", err)
		}
		if !ok {
			return nil, nil, nil, nil, nil, errors.New("access denied to this candidate")
		}
	}

	score, factors, err := u.candidateRepo.GetScoreByCandidateID(ctx, id)
	if err != nil {
		// If score is not found, it's not a fatal error for basic details
		u.log.Debug("candidate score not found", zap.String("candidate_id", id), zap.Error(err))
	}

	// Populate Signed URL for Resume if exists
	if cand.ResumeFileKey != nil && *cand.ResumeFileKey != "" {
		url, err := u.storage.GetFileURL(ctx, *cand.ResumeFileKey)
		if err != nil {
			u.log.Warn("failed to generate signed resume url", zap.String("candidate_id", id), zap.Error(err))
		} else {
			cand.ResumeURL = &url
		}
	}

	return cand, profile, stageID, score, factors, nil
}

func (u *candidateUseCase) GetCandidatesByJob(ctx context.Context, jobID, userID, role string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error) {
	if role != "owner" && role != "admin" {
		ok, err := u.accessRepo.HasAccess(ctx, userID, jobID)
		if err != nil {
			return nil, fmt.Errorf("check access: %w", err)
		}
		if !ok {
			return nil, errors.New("access denied to this job")
		}
	}

	candidatesDTO, err := u.candidateRepo.GetByJobID(ctx, jobID, offset, limit, filter)
	if err != nil {
		return nil, fmt.Errorf("get candidates by job: %w", err)
	}

	return candidatesDTO, nil
}

func (u *candidateUseCase) DeleteCandidate(ctx context.Context, id, actorID, teamID, role string) error {
	if _, _, _, _, _, err := u.GetCandidateByID(ctx, id, actorID, role); err != nil {
		return err
	}

	if err := u.candidateRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete candidate: %w", err)
	}

	_ = u.auditor.Log(ctx, audit.Entry{
		TeamID:    teamID,
		ActorType: audit.ActorUser,
		ActorID:   &actorID,
		Action:    audit.HiringCandidateDeleted,
		TargetID:  &id,
	})

	return nil
}

func (u *candidateUseCase) BulkMoveCandidates(ctx context.Context, p domain.BulkMoveCandidateParams, role string) error {
	if len(p.CandidateIDs) == 0 {
		return nil
	}

	// 1. Validate access to target stage
	targetStage, err := u.pipelineRepo.GetStageByID(ctx, p.ToStageID)
	if err != nil {
		return fmt.Errorf("get target stage: %w", err)
	}

	if role != "owner" && role != "admin" {
		if targetStage.JobID != nil {
			ok, err := u.accessRepo.HasAccess(ctx, p.ChangedBy, *targetStage.JobID)
			if err != nil {
				return fmt.Errorf("check job access: %w", err)
			}
			if !ok {
				return errors.New("access denied to this job")
			}
		}
	}

	// 2. Execute bulk move
	if err := u.candidateRepo.BulkMoveToStage(ctx, p.CandidateIDs, p.ToStageID, p.ChangedBy); err != nil {
		return fmt.Errorf("bulk move repo: %w", err)
	}

	// 3. Audit log (simplified for bulk)
	if targetStage.JobID != nil {
		job, _ := u.jobRepo.GetByID(ctx, *targetStage.JobID)
		if job != nil {
			_ = u.auditor.Log(ctx, audit.Entry{
				TeamID:    job.TeamID,
				ActorType: audit.ActorUser,
				ActorID:   &p.ChangedBy,
				Action:    audit.HiringCandidateMoved, // Reuse existing action or add new one
				Payload: map[string]interface{}{
					"to_stage_id":   p.ToStageID,
					"candidate_ids": p.CandidateIDs,
					"is_bulk":       true,
				},
			})
		}
	}

	return nil
}

func (u *candidateUseCase) MoveCandidate(ctx context.Context, p domain.MoveCandidateParams, role string) error {
	// Validate access
	if _, _, _, _, _, err := u.GetCandidateByID(ctx, p.CandidateID, p.ChangedBy, role); err != nil {
		return err
	}

	// Validate stage transition: only +1 position or terminal stage
	targetStage, err := u.pipelineRepo.GetStageByID(ctx, p.ToStageID)
	if err != nil {
		return fmt.Errorf("get target stage: %w", err)
	}

	// If it's a terminal stage, always allow
	if !targetStage.IsTerminal {
		// Get candidate's current stage
		_, _, currentStageID, err := u.candidateRepo.GetByID(ctx, p.CandidateID)
		if err != nil {
			return fmt.Errorf("get candidate: %w", err)
		}

		if currentStageID != nil {
			currentStage, err := u.pipelineRepo.GetStageByID(ctx, *currentStageID)
			if err != nil {
				return fmt.Errorf("get current stage: %w", err)
			}

			if targetStage.Position != currentStage.Position+1 {
				return ErrInvalidStageTransition
			}
		}
	}

	if err := u.candidateRepo.MoveStage(ctx, p); err != nil {
		return fmt.Errorf("move candidate: %w", err)
	}

	// We should fetch TeamID from targetStage.JobID -> job.TeamID, or from Candidate -> job.TeamID.
	if targetStage.JobID != nil {
		job, err := u.jobRepo.GetByID(ctx, *targetStage.JobID)
		if err != nil {
			u.log.Warn("failed to get job to log candidate_moved audit", zap.Error(err))
		} else {
			_ = u.auditor.Log(ctx, audit.Entry{
				TeamID:    job.TeamID,
				ActorType: audit.ActorUser,
				ActorID:   &p.ChangedBy,
				Action:    audit.HiringCandidateMoved,
				TargetID:  &p.CandidateID,
				Payload:   map[string]string{"to_stage_id": p.ToStageID},
			})
		}
	}

	return nil
}

func (u *candidateUseCase) GetStageHistory(ctx context.Context, candidateID, userID, role string) ([]domain.StageHistoryEntry, error) {
	if _, _, _, _, _, err := u.GetCandidateByID(ctx, candidateID, userID, role); err != nil {
		return nil, err
	}

	entries, err := u.candidateRepo.GetStageHistory(ctx, candidateID)
	if err != nil {
		return nil, fmt.Errorf("get stage history: %w", err)
	}

	return entries, nil
}

func NewCandidateUseCase(log *zap.Logger, candidateRepo repo.CandidateRepository, pipelineRepo repo.PipelineRepository, jobRepo repo.JobRepository, accessRepo repo.AccessRepository, storage storage.FileStorage, publisher *mq.MQPublisher, auditor *audit.Logger) CandidateUseCase {
	return &candidateUseCase{
		log:           log,
		candidateRepo: candidateRepo,
		pipelineRepo:  pipelineRepo,
		jobRepo:       jobRepo,
		accessRepo:    accessRepo,
		storage:       storage,
		publisher:     publisher,
		auditor:       auditor,
	}
}

