package llm

import (
	"backend/internal/domain"
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ─── System prompt ────────────────────────────────────────────────────────────

const chatSystemPrompt = `You are a senior technical recruiter and HR assistant. Your goal is to help users analyze candidate resumes and answer questions based on the provided context.

Context consists of fragments from candidate resumes retrieved via vector search.

Guidelines:
1. Be professional, helpful, and concise.
2. Base your answers strictly on the provided context. If the information is not in the context, say so gracefully.
3. If multiple resumes are mentioned, clarify which one you are referring to.
4. Maintain the conversation flow using the provided message history.
5. Never make up facts or technical skills that are not present in the resumes.
6. Speak the language requested by the user or the provided locale.`

// ─── ChatAssistant ────────────────────────────────────────────────────────────

type ChatAssistant struct {
	provider Provider
}

func NewChatAssistant(provider Provider) *ChatAssistant {
	return &ChatAssistant{
		provider: provider,
	}
}

func (a *ChatAssistant) GenerateResponse(ctx context.Context, teamID string, history []domain.ChatMessage, contextChunks []string, question string, locale string) (string, error) {
	client, settings, err := a.provider.GetClient(ctx, teamID)
	if err != nil {
		return "", fmt.Errorf("failed to get llm client: %w", err)
	}

	cfg := a.provider.GetGlobalConfig()
	maxTokens := cfg.MaxTokensChat
	if maxTokens == 0 {
		maxTokens = 2048
	}

	model := "anthropic/claude-3-sonnet" // Better for chat than haiku
	if settings.ChatModel != nil && *settings.ChatModel != "" {
		model = *settings.ChatModel
	} else if cfg.ChatModel != "" {
		model = cfg.ChatModel
	}

	userContent := buildChatUserPrompt(contextChunks, question, locale)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: chatSystemPrompt,
		},
	}

	// Add history
	for _, msg := range history {
		role := openai.ChatMessageRoleUser
		if msg.Role == domain.RoleAssistant {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add current question
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userContent,
	})

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from llm")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func buildChatUserPrompt(chunks []string, question string, locale string) string {
	var sb strings.Builder

	if locale != "" {
		sb.WriteString(fmt.Sprintf("IMPORTANT: Answer strictly in the '%s' language.\n\n", locale))
	}

	sb.WriteString("CONTEXT FROM RESUMES:\n")
	if len(chunks) > 0 {
		for i, chunk := range chunks {
			sb.WriteString(fmt.Sprintf("--- Fragment %d ---\n%s\n", i+1, chunk))
		}
	} else {
		sb.WriteString("(No relevant resume fragments found)\n")
	}

	sb.WriteString("\nUSER QUESTION:\n")
	sb.WriteString(question)

	return sb.String()
}
