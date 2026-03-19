package activity

import (
	"backend/internal/domain"
	"context"
	"fmt"
)

func (a *Activities) ChatEmbedQuestion(ctx context.Context, input ChatInput) ([]float32, error) {
	// We use the first chunk's embedding as the representation for the whole question
	chunks, err := a.embedder.EmbedChunks(ctx, input.Question, input.TeamID)
	if err != nil {
		return nil, fmt.Errorf("embed question error: %w", err)
	}

	if len(chunks) == 0 {
		return nil, fmt.Errorf("no embedding generated for question")
	}

	return chunks[0].Embedding, nil
}

func (a *Activities) ChatSearchContext(ctx context.Context, input struct {
	ChatInput
	Vector []float32
}) ([]string, error) {
	embeddings, err := a.candidateDB.SearchEmbeddings(ctx, input.TeamID, input.CandidateID, input.Vector, 5)
	if err != nil {
		return nil, fmt.Errorf("search embeddings error: %w", err)
	}

	var chunks []string
	for _, e := range embeddings {
		chunks = append(chunks, e.ChunkText)
	}

	return chunks, nil
}

func (a *Activities) ChatGenerateAnswer(ctx context.Context, input struct {
	ChatInput
	Context []string
}) (ChatOutput, error) {
	answer, err := a.chatAsst.GenerateResponse(ctx, input.TeamID, input.History, input.Context, input.Question, input.Locale)
	if err != nil {
		return ChatOutput{}, fmt.Errorf("generate answer error: %w", err)
	}

	return ChatOutput{Answer: answer}, nil
}

func (a *Activities) ChatSaveMessage(ctx context.Context, msg domain.ChatMessage) error {
	if err := a.chatDB.AddMessage(ctx, &msg); err != nil {
		return fmt.Errorf("save chat message error: %w", err)
	}
	return nil
}
