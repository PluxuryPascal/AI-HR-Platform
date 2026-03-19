package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ─── System prompt ────────────────────────────────────────────────────────────

const interviewSystemPrompt = `You are an expert technical interviewer. Based on the candidate's resume and the desired role, generate a list of 5-8 relevant interview questions and their expected/model answers.

The questions should be a mix of:
1. Technical skills verification (based on the resume).
2. Behavioral questions relevant to the experience.
3. Role-specific scenarios.

Return ONLY a valid JSON object — no markdown, no code blocks, no explanation.

JSON schema:
{
  "interviews": [
    {
      "question": "string",
      "answer":   "string"
    }
  ]
}

Rules:
- Be specific to the candidate's experience.
- The "answer" should be a concise summary of what a "good" response looks like.
- Ensure the language matches the requested locale.`

// ─── JSON response struct ─────────────────────────────────────────────────────

type interviewResponseJSON struct {
	Interviews []interviewPairJSON `json:"interviews"`
}

type interviewPairJSON struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type InterviewGenerator struct {
	provider Provider
}

func NewInterviewGenerator(provider Provider) *InterviewGenerator {
	return &InterviewGenerator{
		provider: provider,
	}
}

type InterviewQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (g *InterviewGenerator) Generate(ctx context.Context, teamID string, resumeText string, locale string) ([]InterviewQuestion, error) {
	client, settings, err := g.provider.GetClient(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm client: %w", err)
	}

	userContent := buildInterviewUserPrompt(resumeText, locale)

	cfg := g.provider.GetGlobalConfig()
	maxTokens := cfg.MaxTokensChat
	if maxTokens == 0 {
		maxTokens = 2048
	}

	model := "anthropic/claude-3-5-sonnet"
	if settings.ChatModel != nil && *settings.ChatModel != "" {
		model = *settings.ChatModel
	} else if cfg.ChatModel != "" {
		model = cfg.ChatModel
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: interviewSystemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userContent,
			},
		},
		MaxTokens:   maxTokens,
		Temperature: 0.5,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate interview questions: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from llm")
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)
	var parsed interviewResponseJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse interview response: %w (raw: %.100s)", err, raw)
	}

	result := make([]InterviewQuestion, 0, len(parsed.Interviews))
	for _, p := range parsed.Interviews {
		result = append(result, InterviewQuestion{
			Question: p.Question,
			Answer:   p.Answer,
		})
	}

	return result, nil
}

func buildInterviewUserPrompt(resumeText string, locale string) string {
	var sb strings.Builder

	if locale != "" {
		sb.WriteString(fmt.Sprintf("IMPORTANT: Generate all questions and answers in the '%s' locale/language.\n\n", locale))
	}

	sb.WriteString("CANDIDATE RESUME:\n")
	sb.WriteString(resumeText)

	return sb.String()
}
