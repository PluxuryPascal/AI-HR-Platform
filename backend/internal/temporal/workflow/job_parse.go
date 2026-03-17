package workflow

import (
	"backend/internal/temporal/activity"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func JobParseWorkflow(ctx workflow.Context, input activity.JobParseInput) (*activity.JobParseOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("JobParseWorkflow started",
		"text_length", len(input.RawText),
		"locale", input.Locale,
	)

	var a activity.Activities
	var output activity.JobParseOutput

	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, jobParseLLMOptions()),
		a.ParseJob,
		input,
	).Get(ctx, &output); err != nil {
		return nil, fmt.Errorf("parse job activity: %w", err)
	}

	logger.Info("JobParseWorkflow completed", "title", output.Title)

	return &output, nil
}

func jobParseLLMOptions() workflow.ActivityOptions {
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
