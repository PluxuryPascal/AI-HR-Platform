package usecase

import (
	"backend/internal/temporal"
	"backend/internal/repo"
	"backend/internal/temporal/activity"
	"backend/internal/temporal/workflow"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type InterviewUseCase interface {
	GenerateQuestions(ctx context.Context, candidateID, teamID, locale string) (activity.InterviewOutput, error)
	GetQuestions(ctx context.Context, candidateID string) (activity.InterviewOutput, error)
}

type interviewUseCase struct {
	log            *zap.Logger
	temporalClient *temporal.Client
	candidateRepo  repo.CandidateRepository
}

func NewInterviewUseCase(log *zap.Logger, temporalClient *temporal.Client, candidateRepo repo.CandidateRepository) InterviewUseCase {
	return &interviewUseCase{
		log:            log,
		temporalClient: temporalClient,
		candidateRepo:  candidateRepo,
	}
}

func (u *interviewUseCase) GenerateQuestions(ctx context.Context, candidateID, teamID, locale string) (activity.InterviewOutput, error) {
	u.log.Info("GenerateQuestions usecase started", zap.String("candidate_id", candidateID))

	input := activity.InterviewInput{
		CandidateID: candidateID,
		TeamID:      teamID,
		Locale:      locale,
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("interview-questions-%s-%d", candidateID, time.Now().Unix()),
		TaskQueue: u.temporalClient.TaskQueue(),
	}

	run, err := u.temporalClient.TemporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.InterviewQuestionsWorkflow, input)
	if err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("failed to execute interview workflow: %w", err)
	}

	var result activity.InterviewOutput
	if err := run.Get(ctx, &result); err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("interview workflow failed: %w", err)
	}

	return result, nil
}

func (u *interviewUseCase) GetQuestions(ctx context.Context, candidateID string) (activity.InterviewOutput, error) {
	_, profile, _, err := u.candidateRepo.GetByID(ctx, candidateID)
	if err != nil {
		return activity.InterviewOutput{}, err
	}

	if profile.InterviewGuide == nil {
		return activity.InterviewOutput{}, fmt.Errorf("interview guide not found")
	}

	var questions []activity.InterviewPair
	if err := json.Unmarshal(profile.InterviewGuide, &questions); err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("failed to unmarshal interview guide: %w", err)
	}

	return activity.InterviewOutput{Questions: questions}, nil
}
