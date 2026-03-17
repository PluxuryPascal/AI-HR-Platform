package workflow

import (
	"backend/internal/temporal/activity"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func EmailGenerateWorkflow(ctx workflow.Context, input activity.EmailGenerateInput) (*activity.EmailGenerateOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("EmailGenerateWorkflow started",
		"candidate_id", input.CandidateID,
		"type", input.EmailType,
		"locale", input.Locale,
	)

	var a activity.Activities
	var output activity.EmailGenerateOutput

	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, emailLLMOptions()),
		a.GenerateEmail,
		input,
	).Get(ctx, &output); err != nil {
		return nil, fmt.Errorf("generate email activity: %w", err)
	}

	logger.Info("EmailGenerateWorkflow completed", "candidate_id", input.CandidateID)

	return &output, nil
}

func emailLLMOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    3,
			InitialInterval:    3 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    15 * time.Second,
		},
	}
}
