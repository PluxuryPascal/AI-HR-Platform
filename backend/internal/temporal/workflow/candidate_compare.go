package workflow

import (
	"backend/internal/temporal/activity"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CompareCandidatesWorkflow(ctx workflow.Context, input activity.CandidateCompareInput) (*activity.CandidateCompareOutput, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("CompareCandidatesWorkflow started")

	var a activity.Activities
	var output activity.CandidateCompareOutput

	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, compareLLMOptions()),
		a.CompareCandidates,
		input,
	).Get(ctx, &output); err != nil {
		return nil, fmt.Errorf("compare candidates activity: %w", err)
	}

	logger.Info("CompareCandidatesWorkflow completed")

	return &output, nil
}

func compareLLMOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    3,
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
		},
	}
}
