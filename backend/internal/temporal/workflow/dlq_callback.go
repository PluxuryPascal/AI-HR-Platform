package workflow

import (
	"backend/internal/temporal/activity"
	"fmt"
	"time"

	hiringv1 "backend/internal/proto/hiring/v1"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

// DLQCallbackWorkflow is a minimal workflow that only calls GRPCCallback
// with FAILED status. Used by the DLQ consumer when all retries are exhausted.
func DLQCallbackWorkflow(ctx workflow.Context, input activity.GRPCCallbackInput) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("DLQCallbackWorkflow started",
		"candidate_id", input.CandidateID,
	)

	input.Status = hiringv1.ParsingStatus_PARSING_STATUS_FAILED

	actCtx := workflow.WithActivityOptions(ctx, dlqCallbackOptions())
	if err := workflow.ExecuteActivity(actCtx, (*activity.Activities).GRPCCallback, input).Get(actCtx, nil); err != nil {
		logger.Error("DLQ callback failed", zap.Error(err))
		return fmt.Errorf("dlq grpc callback: %w", err)
	}

	logger.Info("DLQCallbackWorkflow completed", "candidate_id", input.CandidateID)

	return nil
}

func dlqCallbackOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 15 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    5,
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 1.5,
			MaximumInterval:    10 * time.Second,
		},
	}
}
