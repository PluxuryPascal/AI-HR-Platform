package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type CandidateCompareInput struct {
	ID         string
	FirstName  string
	LastName   string
	Role       string
	MatchScore int
	Skills     []string
	ParsedText string
}

type CandidateCompareResult struct {
	Experience  string   `json:"experience"`
	Skills      []string `json:"skills"`
	SalaryRange string   `json:"salary_range"`
	Risks       string   `json:"risks"`
	Summary     string   `json:"summary"`
}

type CompareResult map[string]CandidateCompareResult

// ─── System prompt ────────────────────────────────────────────────────────────

const comparatorSystemPrompt = `You are a senior technical recruiter performing a side-by-side candidate comparison.

Analyze ALL candidates provided and return a JSON object where each key is a candidate ID and the value is their comparison profile.

Return ONLY a valid JSON object — no markdown, no code blocks, no explanation.

Output schema (one entry per candidate ID):
{
  "<candidate_id>": {
    "experience":   string,   // e.g. "5 years", "3+ years", "< 1 year"
    "skills":       string[], // top 4-6 most relevant skills for the role
    "salary_range": string,   // estimated range e.g. "$90k - $110k", or "" if unknown
    "risks":        string,   // main hiring risk, or "None identified"
    "summary":      string    // 1-2 sentences on fit for the role
  }
}

Rules:
- Base experience and skills on the resume text — do not fabricate.
- salary_range: estimate based on role title, experience level, and market data. Use USD unless location suggests otherwise. Empty string if too uncertain.
- risks: be specific and honest — e.g. "High salary expectations", "Remote only", "Overqualified for junior role", "Gap in employment 2022-2023". Use "None identified" if genuinely no concerns.
- summary: compare the candidate to the role requirements, not to other candidates. Focus on fit.
- Return ALL candidate IDs provided — never omit one.`

type CandidateComparator struct {
	provider Provider
}

func NewCandidateComparator(provider Provider) *CandidateComparator {
	return &CandidateComparator{provider: provider}
}

func (c *CandidateComparator) Compare(
	ctx context.Context,
	candidates []CandidateCompareInput,
	jobRequirements, locale string,
	teamID string,
) (CompareResult, error) {
	if len(candidates) < 2 {
		return nil, fmt.Errorf("comparison requires at least 2 candidates, got %d", len(candidates))
	}

	client, settings, err := c.provider.GetClient(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm client: %w", err)
	}

	cfg := c.provider.GetGlobalConfig()
	model := "anthropic/claude-3-5-sonnet"
	if settings.ChatModel != nil && *settings.ChatModel != "" {
		model = *settings.ChatModel
	} else if cfg.ChatModel != "" {
		model = cfg.ChatModel
	}

	maxTokens := cfg.MaxTokensChat
	if maxTokens == 0 {
		maxTokens = 2048
	}

	maxTokens = max(maxTokens, len(candidates)*300)

	userContent := buildCompareUserPrompt(candidates, jobRequirements, locale)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: comparatorSystemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userContent},
		},
		MaxTokens:   maxTokens,
		Temperature: 0.3,
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

	var result CompareResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("parse compare json: %w (raw: %.300s)", err, raw)
	}

	for _, candidate := range candidates {
		if _, ok := result[candidate.ID]; !ok {
			result[candidate.ID] = CandidateCompareResult{
				Experience:  "Unknown",
				Skills:      candidate.Skills,
				SalaryRange: "",
				Risks:       "None identified",
				Summary:     "Analysis unavailable.",
			}
		}
	}

	return result, nil
}

func buildCompareUserPrompt(candidates []CandidateCompareInput, jobRequirements, locale string) string {
	var sb strings.Builder

	if locale != "" {
		sb.WriteString(fmt.Sprintf("IMPORTANT: You must write all output text, including experience details, risks, and summary, in the '%s' locale/language.\n\n", locale))
	}

	if jobRequirements != "" {
		sb.WriteString("Job requirements:\n")
		sb.WriteString(jobRequirements)
		sb.WriteString("\n\n")
	}

	sb.WriteString("Candidates to compare:\n\n")

	for _, c := range candidates {
		sb.WriteString(fmt.Sprintf("--- Candidate ID: %s ---\n", c.ID))
		sb.WriteString(fmt.Sprintf("Name:        %s %s\n", c.FirstName, c.LastName))
		sb.WriteString(fmt.Sprintf("Applied for: %s\n", c.Role))
		sb.WriteString(fmt.Sprintf("Match score: %d%%\n", c.MatchScore))

		if len(c.Skills) > 0 {
			skills := c.Skills
			if len(skills) > 8 {
				skills = skills[:8]
			}
			sb.WriteString(fmt.Sprintf("Known skills: %s\n", strings.Join(skills, ", ")))
		}

		if c.ParsedText != "" {
			resumeExcerpt := c.ParsedText
			if len(resumeExcerpt) > 2000 {
				resumeExcerpt = resumeExcerpt[:2000] + "…"
			}
			sb.WriteString("Resume excerpt:\n")
			sb.WriteString(resumeExcerpt)
		}

		sb.WriteString("\n\n")
	}

	return sb.String()
}
