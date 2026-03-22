package activity

import (
	"backend/internal/llm"
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
)

func (a *Activities) CompareCandidates(ctx context.Context, input CandidateCompareInput) (*CandidateCompareOutput, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Comparing candidates started", "candidates_count", len(input.Candidates), "job_requirements", input.JobRequirements)

	candidates := make([]llm.CandidateCompareInput, len(input.Candidates))
	for i, candidate := range input.Candidates {
		candidates[i] = llm.CandidateCompareInput{
			ID:         candidate.ID,
			FirstName:  candidate.FirstName,
			LastName:   candidate.LastName,
			Role:       candidate.Role,
			MatchScore: candidate.MatchScore,
			Skills:     candidate.Skills,
			ParsedText: candidate.ParsedText,
		}
	}

	result, err := a.comparator.Compare(ctx, candidates, input.JobRequirements, input.Locale, input.TeamID)
	if err != nil {
		logger.Error("Failed to compare candidates", "error", err)
		return nil, fmt.Errorf("failed to compare candidates: %w", err)
	}

	output := make(CandidateCompareOutput, len(result))
	for i, candidate := range result {
		output[i] = CompareResultEntry{
			Experience:  candidate.Experience,
			Skills:      candidate.Skills,
			SalaryRange: candidate.SalaryRange,
			Risks:       candidate.Risks,
			Summary:     candidate.Summary,
		}
	}

	logger.Info("Comparing candidates completed", "candidates_count", len(input.Candidates))

	return &output, nil
}
