package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ResumePipelineWorkflow(ctx workflow.Context, input []byte) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var parsedText string
	if err := workflow.ExecuteActivity(ctx, "resume-pipeline-activity", input).Get(ctx, &parsedText); err != nil {
		return err
	}

	llmAO := workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, llmAO)

	var llmResponse string
	if err := workflow.ExecuteActivity(ctx, "llm-activity", parsedText).Get(ctx, &llmResponse); err != nil {
		return err
	}

	return nil
}
