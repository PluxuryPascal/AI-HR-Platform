package workflow

import (
	"backend/internal/domain"
	"backend/internal/temporal/activity"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

func ChatWorkflow(ctx workflow.Context, input activity.ChatInput) (activity.ChatOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ChatWorkflow started", zap.String("session_id", input.SessionID))

	var a *activity.Activities

	// 1. Embed Question
	var vector []float32
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, chatActivityOptions()),
		a.ChatEmbedQuestion,
		input,
	).Get(ctx, &vector); err != nil {
		return activity.ChatOutput{}, err
	}

	// 2. Search Context
	var chunks []string
	searchParams := struct {
		activity.ChatInput
		Vector []float32
	}{
		ChatInput: input,
		Vector:    vector,
	}
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, chatActivityOptions()),
		a.ChatSearchContext,
		searchParams,
	).Get(ctx, &chunks); err != nil {
		return activity.ChatOutput{}, err
	}

	// 3. Generate Answer
	var output activity.ChatOutput
	genParams := struct {
		activity.ChatInput
		Context []string
	}{
		ChatInput: input,
		Context:   chunks,
	}
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, chatActivityOptions()),
		a.ChatGenerateAnswer,
		genParams,
	).Get(ctx, &output); err != nil {
		return activity.ChatOutput{}, err
	}

	// 4. Save AI Response
	aiMsg := domain.ChatMessage{
		SessionID: input.SessionID,
		Role:      domain.RoleAssistant,
		Content:   output.Answer,
	}
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, chatSaveActivityOptions()),
		a.ChatSaveMessage,
		aiMsg,
	).Get(ctx, nil); err != nil {
		logger.Error("failed to save AI chat message", zap.Error(err))
	}

	logger.Info("ChatWorkflow completed", zap.String("session_id", input.SessionID))
	return output, nil
}

func chatActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    3,
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 2.0,
		},
	}
}

func chatSaveActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout: 15 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    5,
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 1.5,
		},
	}
}
