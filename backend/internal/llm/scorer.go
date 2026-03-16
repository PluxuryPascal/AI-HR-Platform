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

const scorerSystemPrompt = `You are a senior technical recruiter. Score how well a candidate's resume matches the given job requirements.
 
Return ONLY a valid JSON object — no markdown, no code blocks, no explanation.
 
JSON schema:
{
  "match_score": integer (0-100),
  "factors": [
    {
      "type":        "positive" | "negative" | "neutral",
      "description": string,
      "impact":      integer (1-100)
    }
  ]
}
 
Scoring rules:
- "match_score": overall match percentage 0-100. 
  90-100: exceptional match. 70-89: strong match. 50-69: partial match. Below 50: weak match.
- "factors": list 3-7 most significant factors influencing the score.
  - "positive": candidate exceeds or meets a key requirement.
  - "negative": candidate lacks a required skill or has insufficient experience.
  - "neutral":  candidate has relevant experience not in the requirements (nice-to-have or risk).
  - "impact": how much this factor influenced the score (1-100). All impacts should sum to ~100.
- Be specific: mention concrete skills, years of experience, technologies.
- If job requirements are empty, score based on general employability.`

// ─── JSON response struct ─────────────────────────────────────────────────────

type scoreResponseJSON struct {
	MatchScore int               `json:"match_score"`
	Factors    []scoreFactorJSON `json:"factors"`
}

type scoreFactorJSON struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Impact      int    `json:"impact"`
}

type Scorer struct {
	client *openai.Client
	cfg    *config.OpenRouter
}

func NewScorer(client *openai.Client, cfg *config.OpenRouter) *Scorer {
	return &Scorer{
		client: client,
		cfg:    cfg,
	}
}

func (s *Scorer) Score(ctx context.Context, resumeText, jobRequirements string) (*domain.ScoreResult, error) {
	userContent := buildScoreUserPrompt(resumeText, jobRequirements)

	maxTokens := s.cfg.MaxTokensScore
	if maxTokens == 0 {
		maxTokens = 1024
	}

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.cfg.ScoreModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: scorerSystemPrompt,
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

	var parsed scoreResponseJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse scorer response: %w", err)
	}

	result := toDomainScoreResult(parsed)

	return result, nil
}

func buildScoreUserPrompt(resumeText, jobRequirements string) string {
	var sb strings.Builder

	sb.WriteString("Job requirements:\n")
	if jobRequirements != "" {
		sb.WriteString(jobRequirements)
	} else {
		sb.WriteString("(not provided — score general employability)")
	}

	sb.WriteString("\n\nCandidate resume:\n")
	sb.WriteString(resumeText)

	return sb.String()
}

func toDomainScoreResult(r scoreResponseJSON) *domain.ScoreResult {
	// Clamp score to valid range.
	score := r.MatchScore
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	factors := make([]domain.ScoreFactor, 0, len(r.Factors))
	for _, f := range r.Factors {
		ft := domain.FactorType(f.Type)
		// Guard against unexpected type values from LLM.
		if ft != domain.FactorPositive && ft != domain.FactorNegative && ft != domain.FactorNeutral {
			ft = domain.FactorNeutral
		}

		impact := f.Impact
		if impact < 1 {
			impact = 1
		}
		if impact > 100 {
			impact = 100
		}

		factors = append(factors, domain.ScoreFactor{
			Type:        ft,
			Description: f.Description,
			Impact:      impact,
		})
	}

	return &domain.ScoreResult{
		MatchScore: score,
		Factors:    factors,
	}
}
