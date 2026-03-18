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
	GetCandidateByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error)
	GetCandidatesByJob(ctx context.Context, jobID string, offset, limit int, filter domain.CandidateFilter) (*domain.CandidatesDTO, error)
	DeleteCandidate(ctx context.Context, id, actorID, teamID string) error
	MoveCandidate(ctx context.Context, p domain.MoveCandidateParams) error
	FinalizeAIParsing(ctx context.Context, result domain.AIParsingResult) error
	UpdateByRecruiter(ctx context.Context, candidate *domain.Candidate) error
	ConfirmManualReview(ctx context.Context, candidate *domain.Candidate) error
	GetStageHistory(ctx context.Context, candidateID string) ([]domain.StageHistoryEntry, error)
}

type candidateUseCase struct {
	log           *zap.Logger
	candidateRepo repo.CandidateRepository
	pipelineRepo  repo.PipelineRepository
	jobRepo       repo.JobRepository
	storage       storage.FileStorage
	publisher     *mq.MQPublisher
	auditor       *audit.Logger
}

func (u *candidateUseCase) ConfirmManualReview(ctx context.Context, candidate *domain.Candidate) error {
	current, _, _, err := u.candidateRepo.GetByID(ctx, candidate.ID)
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

func (u *candidateUseCase) UpdateByRecruiter(ctx context.Context, candidate *domain.Candidate) error {
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

func (u *candidateUseCase) DeleteCandidate(ctx context.Context, id, actorID, teamID string) error {
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

func (u *candidateUseCase) MoveCandidate(ctx context.Context, p domain.MoveCandidateParams) error {
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

func (u *candidateUseCase) GetStageHistory(ctx context.Context, candidateID string) ([]domain.StageHistoryEntry, error) {
	entries, err := u.candidateRepo.GetStageHistory(ctx, candidateID)
	if err != nil {
		return nil, fmt.Errorf("get stage history: %w", err)
	}

	return entries, nil
}

func NewCandidateUseCase(log *zap.Logger, candidateRepo repo.CandidateRepository, pipelineRepo repo.PipelineRepository, jobRepo repo.JobRepository, storage storage.FileStorage, publisher *mq.MQPublisher, auditor *audit.Logger) CandidateUseCase {
	return &candidateUseCase{
		log:           log,
		candidateRepo: candidateRepo,
		pipelineRepo:  pipelineRepo,
		jobRepo:       jobRepo,
		storage:       storage,
		publisher:     publisher,
		auditor:       auditor,
	}
}

