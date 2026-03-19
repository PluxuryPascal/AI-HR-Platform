package llm

import (
	"backend/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type EmailTone string

const (
	EmailToneProfessional EmailTone = "professional"
	EmailToneFriendly     EmailTone = "friendly"
	EmailToneBrief        EmailTone = "brief"
)

type EmailGenerateInput struct {
	TeamID             string
	CandidateFirstName string
	CandidateLastName  string
	CandidateEmail     string
	Role               string
	Skills             []string
	MatchScore         int

	Type   domain.EmailType
	Tone   EmailTone
	Locale string

	RecruiterName string
	CompanyName   string
}

type EmailGenerateResult struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// ─── System prompt ────────────────────────────────────────────────────────────

const emailGeneratorSystemPrompt = `You are an expert HR communications specialist. Write a personalized recruitment email based on the provided candidate data and parameters.

Return ONLY a valid JSON object — no markdown, no code blocks, no explanation.

JSON schema:
{
  "subject": string,
  "body":    string
}

Writing rules:
- Write the email entirely in the language specified by the "locale" field (e.g. "ru" → Russian, "es" → Spanish, "en" → English).
- Adapt tone strictly:
  "professional" — formal address (Dear/Уважаемый/Estimado), structured paragraphs, no emoji.
  "friendly"     — first name only, warm conversational tone, 1-2 emoji allowed.
  "brief"        — 3-5 sentences max, no pleasantries, direct.
- For "rejection": acknowledge one genuine strength from their skills, give a vague but honest reason, wish them well. Never make up a reason — use only the data provided.
- For "interview_invite": highlight the 1-2 most relevant skills, express genuine enthusiasm, ask for availability.
- Body must use \n for line breaks. Do NOT use HTML.
- If RecruiterName is provided, sign with it. Otherwise sign as "Hiring Team".
- If CompanyName is provided, mention it naturally.`

type EmailGenerator struct {
	provider Provider
}

func NewEmailGenerator(provider Provider) *EmailGenerator {
	return &EmailGenerator{provider: provider}
}

func (g *EmailGenerator) Generate(ctx context.Context, input EmailGenerateInput) (*EmailGenerateResult, error) {
	client, settings, err := g.provider.GetClient(ctx, input.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm client: %w", err)
	}

	cfg := g.provider.GetGlobalConfig()
	model := "anthropic/claude-3-haiku"
	if settings.ChatModel != nil && *settings.ChatModel != "" {
		model = *settings.ChatModel
	} else if cfg.ChatModel != "" {
		model = cfg.ChatModel
	}

	maxTokens := cfg.MaxTokensChat
	if maxTokens == 0 {
		maxTokens = 2048
	}

	userContent := buildEmailUserPrompt(input)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: emailGeneratorSystemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userContent},
		},
		MaxTokens:   maxTokens,
		Temperature: 0.7,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("openrouter chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openrouter returned empty choices")
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)

	var result EmailGenerateResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("parse email json: %w (raw: %.200s)", err, raw)
	}

	if result.Subject == "" || result.Body == "" {
		return nil, fmt.Errorf("llm returned empty subject or body")
	}

	return &result, nil
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func buildEmailUserPrompt(input EmailGenerateInput) string {
	var sb strings.Builder

	emailTypeStr := "rejection"
	if input.Type == domain.EmailInterviewInvite {
		emailTypeStr = "interview_invite"
	}

	sb.WriteString(fmt.Sprintf("locale: %s\n", localeToLanguage(input.Locale)))
	sb.WriteString(fmt.Sprintf("email_type: %s\n", emailTypeStr))
	sb.WriteString(fmt.Sprintf("tone: %s\n", string(input.Tone)))
	sb.WriteString("\nCandidate:\n")
	sb.WriteString(fmt.Sprintf("  first_name: %s\n", input.CandidateFirstName))
	sb.WriteString(fmt.Sprintf("  last_name:  %s\n", input.CandidateLastName))
	sb.WriteString(fmt.Sprintf("  role:       %s\n", input.Role))
	sb.WriteString(fmt.Sprintf("  match_score: %d%%\n", input.MatchScore))

	if len(input.Skills) > 0 {
		skills := input.Skills
		if len(skills) > 5 {
			skills = skills[:5]
		}
		sb.WriteString(fmt.Sprintf("  skills: %s\n", strings.Join(skills, ", ")))
	}

	if input.RecruiterName != "" {
		sb.WriteString(fmt.Sprintf("\nrecruiter_name: %s\n", input.RecruiterName))
	}
	if input.CompanyName != "" {
		sb.WriteString(fmt.Sprintf("company_name: %s\n", input.CompanyName))
	}

	return sb.String()
}

func localeToLanguage(locale string) string {
	if idx := strings.Index(locale, "-"); idx != -1 {
		locale = locale[:idx]
	}

	switch locale {
	case "ru":
		return "Russian"
	case "es":
		return "Spanish"
	default:
		return "English"
	}
}
