package llm

import (
	"backend/internal/domain"
	"backend/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ─── System prompt ────────────────────────────────────────────────────────────

const resumeParserSystemPrompt = `You are an expert resume parser. Extract candidate information from the provided resume text and return ONLY a valid JSON object — no markdown, no code blocks, no explanation.
 
JSON schema (use null for fields you cannot find):
{
  "first_name":       string | null,
  "last_name":        string | null,
  "email":            string | null,
  "location":         string | null,
  "skills":           string[],
  "experience_years": integer | null,
  "education":        string | null,
  "summary":          string | null
}
 
Rules:
- "skills": extract ALL technical and soft skills mentioned. Return as an array of short strings.
- "experience_years": total years of professional experience as an integer. Null if not determinable.
- "education": highest degree and field, e.g. "MSc Computer Science". Null if not found.
- "summary": 2-3 sentence professional summary in third person. Null if resume has no usable content.
- "email": extract email address if present. Null if not found.
- Return null for any field you are not confident about — do not guess.`

// ─────────────────────────────────────────────────────────────────────────────

type parseResponseJSON struct {
	FirstName       *string  `json:"first_name"`
	LastName        *string  `json:"last_name"`
	Email           *string  `json:"email"`
	Location        *string  `json:"location"`
	Skills          []string `json:"skills"`
	ExperienceYears *int     `json:"experience_years"`
	Education       *string  `json:"education"`
	Summary         *string  `json:"summary"`
}

type ResumeParser struct {
	client *openai.Client
	cfg    *config.OpenRouter
}

func NewResumeParser(client *openai.Client, cfg *config.OpenRouter) *ResumeParser {
	return &ResumeParser{
		client: client,
		cfg:    cfg,
	}
}

func (p *ResumeParser) Parse(ctx context.Context, resumeText, jobRequirements, locale string) (*domain.ParseResult, error) {
	userContent := buildParseUserPrompt(resumeText, jobRequirements, locale)

	maxTokens := p.cfg.MaxTokensParse
	if maxTokens == 0 {
		maxTokens = 1024
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: p.cfg.ParseModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: resumeParserSystemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userContent,
			},
		},
		MaxTokens:   maxTokens,
		Temperature: 0,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in chat completion response")
	}

	raw := strings.TrimSpace(resp.Choices[0].Message.Content)

	var parsed parseResponseJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse resume parser response: %w", err)
	}

	result := toDomainParseResult(parsed)
	result.MissingFields = collectMissingFields(result)

	return result, nil
}

func buildParseUserPrompt(resumeText, jobRequirements, locale string) string {
	var sb strings.Builder

	if locale != "" {
		sb.WriteString(fmt.Sprintf("IMPORTANT: You must write the summary and any descriptive strings in the '%s' locale/language.\n\n", locale))
	}

	sb.WriteString("Resume text:\n")
	sb.WriteString(resumeText)

	if jobRequirements != "" {
		sb.WriteString("\n\nJob requirements (use as context for skills extraction and summary):\n")
		sb.WriteString(jobRequirements)
	}

	return sb.String()
}

func toDomainParseResult(r parseResponseJSON) *domain.ParseResult {
	result := &domain.ParseResult{
		FirstName:       r.FirstName,
		LastName:        r.LastName,
		Email:           r.Email,
		Location:        r.Location,
		Skills:          r.Skills,
		ExperienceYears: r.ExperienceYears,
		Education:       r.Education,
		Summary:         r.Summary,
	}

	if result.Skills == nil {
		result.Skills = []string{}
	}

	structured, err := json.Marshal(r)
	if err == nil {
		result.StructuredData = structured
	}

	return result
}

func collectMissingFields(r *domain.ParseResult) []string {
	var missing []string

	if r.FirstName == nil || *r.FirstName == "" {
		missing = append(missing, "first_name")
	}
	if r.LastName == nil || *r.LastName == "" {
		missing = append(missing, "last_name")
	}
	if r.Email == nil || *r.Email == "" {
		missing = append(missing, "email")
	}

	return missing
}
