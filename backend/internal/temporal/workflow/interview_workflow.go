package workflow

import (
	"backend/internal/temporal/activity"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

func InterviewQuestionsWorkflow(ctx workflow.Context, input activity.InterviewInput) (activity.InterviewOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("InterviewQuestionsWorkflow started", "candidate_id", input.CandidateID)

	var a activity.Activities

	var resumeText string
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, interviewActivityOptions()),
		a.InterviewGetCandidateText,
		input,
	).Get(ctx, &resumeText); err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("failed to get candidate text: %w", err)
	}

	var output activity.InterviewOutput
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, interviewLLMActivityOptions()),
		a.InterviewGenerateQuestions,
		input,
		resumeText,
	).Get(ctx, &output); err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("failed to generate interview questions: %w", err)
	}

	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, interviewActivityOptions()),
		a.InterviewSaveGuide,
		input,
		output,
	).Get(ctx, nil); err != nil {
		return activity.InterviewOutput{}, fmt.Errorf("failed to save interview guide: %w", err)
	}

	logger.Info("InterviewQuestionsWorkflow completed", zap.String("candidate_id", input.CandidateID))
	return output, nil
}

func interviewActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    3,
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
		},
	}
}

func interviewLLMActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    3,
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 2.0,
		},
	}
}
