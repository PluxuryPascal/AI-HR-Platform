package usecase

import (
	"backend/internal/temporal/activity"
	"backend/internal/temporal/workflow"
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type InterviewUseCase interface {
	GenerateQuestions(ctx context.Context, candidateID, teamID, locale string) (activity.InterviewOutput, error)
}

type interviewUseCase struct {
	log            *zap.Logger
	temporalClient client.Client
}

func NewInterviewUseCase(log *zap.Logger, temporalClient client.Client) InterviewUseCase {
	return &interviewUseCase{
		log:            log,
		temporalClient: temporalClient,
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
		TaskQueue: "interview-tasks",
	}

	run, err := u.temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.InterviewQuestionsWorkflow, input)
	if err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("failed to execute interview workflow: %w", err)
	}

	var result activity.InterviewOutput
	if err := run.Get(ctx, &result); err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("interview workflow failed: %w", err)
	}

	return result, nil
}
