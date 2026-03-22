package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type JobParseResult struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Requirements []string `json:"requirements"`
	Locale       string   `json:"-"`
	WorkFormat   string   `json:"work_format"`
	SalaryMin    int      `json:"salary_min"`
	SalaryMax    int      `json:"salary_max"`
	Currency     string   `json:"currency"`
}

// ─── System prompt ────────────────────────────────────────────────────────────

const jobParserSystemPrompt = `You are an expert HR analyst. Extract structured job posting information from the provided raw text.

Return ONLY a valid JSON object — no markdown, no code blocks, no explanation.

JSON schema:
{
  "title":        string,
  "description":  string,
  "requirements": string[],
  "work_format":  "Remote" | "Onsite" | "Hybrid",
  "salary_min":   integer,
  "salary_max":   integer,
  "currency":     string
}

Extraction rules:
- "title": extract the exact job title. If multiple titles are mentioned, use the most senior/specific one.
- "description": the full job description body — responsibilities, company context, benefits. Keep it clean but complete. Do not include requirements here.
  IMPORTANT: The description should be detailed and contain at least 10 sentences.
  IMPORTANT: write the description in the language specified by the "output_language" field in the user message, regardless of the input text language.
- "requirements": extract EACH requirement as a separate, atomic string. Be specific:
    GOOD: ["5+ years of Go experience", "Familiarity with Kubernetes", "Upper-Intermediate English"]
    BAD:  ["Strong technical background", "Good communication skills"]
  Include hard requirements (must-have) and soft requirements (nice-to-have) as separate items.
  Aim for 5-15 items. Never return placeholder text.
- "work_format": look for keywords like "remote", "hybrid", "office", "on-site", "в офисе", "удалённо", "гибридный". Default to "Onsite".
- "salary_min" / "salary_max": extract numeric values only. If not explicitly mentioned in the text, you MUST estimate a reasonable market salary range based on the job title and requirements.
- "currency": "USD", "RUB", "EUR", etc. Empty string if not mentioned.`

type JobParser struct {
	provider Provider
}

func NewJobParser(provider Provider) *JobParser {
	return &JobParser{provider: provider}
}

func (p *JobParser) Parse(ctx context.Context, rawText, locale string, teamID string) (*JobParseResult, error) {
	if strings.TrimSpace(rawText) == "" {
		return nil, fmt.Errorf("job parser: input text is empty")
	}

	client, settings, err := p.provider.GetClient(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm client: %w", err)
	}

	cfg := p.provider.GetGlobalConfig()
	model := "anthropic/claude-3-haiku"
	if settings.ParseModel != nil && *settings.ParseModel != "" {
		model = *settings.ParseModel
	} else if cfg.ParseModel != "" {
		model = cfg.ParseModel
	}

	maxTokens := cfg.MaxTokensParse
	if maxTokens == 0 {
		maxTokens = 1024
	}

	if maxTokens < 1500 {
		maxTokens = 1500
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: jobParserSystemPrompt},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: buildJobParserUserPrompt(rawText, locale),
			},
		},
		MaxTokens:   maxTokens,
		Temperature: 0,
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

	var result JobParseResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("parse job json: %w (raw: %.200s)", err, raw)
	}

	sanitise(&result)

	return &result, nil
}

func buildJobParserUserPrompt(rawText, locale string) string {
	var sb strings.Builder
	sb.WriteString("output_language: ")
	sb.WriteString(localeToLanguage(locale))
	sb.WriteString("\n\nJob description text:\n\n")
	sb.WriteString(rawText)
	return sb.String()
}

func sanitise(r *JobParseResult) {
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
	r.Currency = strings.TrimSpace(strings.ToUpper(r.Currency))

	validFormats := map[string]bool{
		"Remote": true, "Onsite": true, "Hybrid": true,
	}
	if !validFormats[r.WorkFormat] {
		r.WorkFormat = "Onsite"
	}

	cleaned := r.Requirements[:0]
	for _, req := range r.Requirements {
		if trimmed := strings.TrimSpace(req); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	r.Requirements = cleaned

	if r.SalaryMin < 0 {
		r.SalaryMin = 0
	}
	if r.SalaryMax < 0 {
		r.SalaryMax = 0
	}
}
