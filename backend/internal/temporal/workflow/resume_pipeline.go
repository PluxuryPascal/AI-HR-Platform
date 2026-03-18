package workflow

import (
	"backend/internal/domain"
	hiringv1 "backend/internal/proto/hiring/v1"
	"backend/internal/temporal/activity"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

func ResumePipelineWorkflow(ctx workflow.Context, input activity.ResumePipelineInput) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("ResumePipelineWorkflow started", "candidate id", input.CandidateID)

	var a activity.Activities

	callBackInput := &activity.GRPCCallbackInput{
		CandidateID: input.CandidateID,
		JobID:       input.JobID,
		Status:      hiringv1.ParsingStatus_PARSING_STATUS_SUCCESS,
	}

	defer func() {
		callBackCtx := workflow.WithActivityOptions(ctx, callbackActivityOptions())
		if err := workflow.ExecuteActivity(callBackCtx, a.GRPCCallback, callBackInput).Get(callBackCtx, nil); err != nil {
			logger.Error("Failed to call back", zap.Error(err))
		}
	}()

	var parsedText string
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, pdfActivityOptions()),
		a.PDFExtract,
		input,
	).Get(ctx, &parsedText); err != nil {
		return fmt.Errorf("failed to extract pdf: %w", err)
	}

	var parseResult *domain.ParseResult
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, llmActivityOptions()),
		a.LLMParse,
		activity.LLMParseInput{
			CandidateID: input.CandidateID,
			JobID:       input.JobID,
			ParsedText:  parsedText,
			Locale:      input.Locale,
		},
	).Get(ctx, &parseResult); err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	callBackInput.ParseResult = *parseResult
	callBackInput.ParsedText = parsedText
	callBackInput.Status = domainStatusToProto(parseResult.ParsingStatus())

	sel := workflow.NewSelector(ctx)

	var scoreErr, embedErr error

	sel.AddFuture(workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, llmActivityOptions()),
		a.LLMScore,
		activity.LLMScoreInput{
			CandidateID:     input.CandidateID,
			ParsedText:      parsedText,
			JobRequirements: parseResult.JobRequirements,
			Locale:          input.Locale,
		},
	), func(f workflow.Future) {
		scoreErr = f.Get(ctx, nil)
	})

	sel.AddFuture(workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, embedActivityOptions()),
		a.LLMEmbed,
		activity.LLMEmbedInput{
			CandidateID: input.CandidateID,
			TeamID:      input.TeamID,
			ParsedText:  parsedText,
		},
	), func(f workflow.Future) {
		embedErr = f.Get(ctx, nil)
	})

	sel.Select(ctx)

	if scoreErr != nil {
		return fmt.Errorf("failed to score: %w", scoreErr)
	}

	if embedErr != nil {
		return fmt.Errorf("failed to embed: %w", embedErr)
	}

	logger.Info("ResumePipelineWorkflow completed", zap.String("candidate id", input.CandidateID))

	return nil
}

func pdfActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        3,
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			NonRetryableErrorTypes: []string{"NonRetryablePDFError"},
		},
	}
}

func llmActivityOptions() workflow.ActivityOptions {
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

func embedActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    3,
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 2.0,
		},
	}
}

func callbackActivityOptions() workflow.ActivityOptions {
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

func domainStatusToProto(s domain.CandidateParsingStatus) hiringv1.ParsingStatus {
	switch s {
	case domain.ParsingStatusNeedsReview:
		return hiringv1.ParsingStatus_PARSING_STATUS_NEEDS_REVIEW
	default:
		return hiringv1.ParsingStatus_PARSING_STATUS_SUCCESS
	}
}
