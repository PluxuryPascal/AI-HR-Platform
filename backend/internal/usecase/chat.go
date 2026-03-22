package usecase

import (
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/internal/temporal"
	"backend/internal/temporal/activity"
	"backend/internal/temporal/workflow"
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type ChatUseCase interface {
	Chat(ctx context.Context, teamID string, userID string, req domain.ChatRequest) (string, string, error)
	GetHistory(ctx context.Context, sessionID string) ([]domain.ChatMessage, error)
	ListSessions(ctx context.Context, teamID, userID string) ([]domain.ChatSession, error)
}

type chatUseCase struct {
	log            *zap.Logger
	chatRepo       repo.ChatRepository
	temporalClient *temporal.Client
}

func NewChatUseCase(log *zap.Logger, chatRepo repo.ChatRepository, temporalClient *temporal.Client) ChatUseCase {
	return &chatUseCase{
		log:            log,
		chatRepo:       chatRepo,
		temporalClient: temporalClient,
	}
}

func (u *chatUseCase) Chat(ctx context.Context, teamID string, userID string, req domain.ChatRequest) (string, string, error) {
	// 1. Find or create session
	chatType := domain.ChatGlobalSearch
	if req.CandidateID != nil {
		chatType = domain.ChatLocalCandidate
	}

	session, err := u.chatRepo.FindOrCreateSession(ctx, teamID, userID, chatType, req.CandidateID)
	if err != nil {
		return "", "", fmt.Errorf("find or create session: %w", err)
	}

	// 2. Save user message
	userMsg := domain.ChatMessage{
		SessionID: session.ID,
		Role:      domain.RoleUser,
		Content:   req.Question,
	}
	if err := u.chatRepo.AddMessage(ctx, &userMsg); err != nil {
		return "", "", fmt.Errorf("save user message: %w", err)
	}

	// 3. Get history for context
	history, err := u.chatRepo.GetHistory(ctx, session.ID)
	if err != nil {
		u.log.Warn("failed to get chat history for context", zap.Error(err))
	}

	// Limit history to last 10 messages for context
	if len(history) > 10 {
		history = history[len(history)-10:]
	}

	// 4. Run Workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("chat-%s", session.ID),
		TaskQueue: u.temporalClient.TaskQueue(),
	}

	input := activity.ChatInput{
		SessionID:   session.ID,
		TeamID:      teamID,
		CandidateID: req.CandidateID,
		Question:    req.Question,
		Locale:      req.Locale,
		History:     history,
	}

	run, err := u.temporalClient.TemporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.ChatWorkflow, input)
	if err != nil {
		return "", "", fmt.Errorf("failed to start chat workflow: %w", err)
	}

	var output activity.ChatOutput
	err = run.Get(ctx, &output)
	if err != nil {
		return "", "", fmt.Errorf("chat workflow failed: %w", err)
	}

	return output.Answer, session.ID, nil
}

func (u *chatUseCase) GetHistory(ctx context.Context, sessionID string) ([]domain.ChatMessage, error) {
	return u.chatRepo.GetHistory(ctx, sessionID)
}

func (u *chatUseCase) ListSessions(ctx context.Context, teamID, userID string) ([]domain.ChatSession, error) {
	return u.chatRepo.GetSessions(ctx, teamID, userID)
}
