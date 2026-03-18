package activity

import (
	"backend/internal/llm"
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.uber.org/zap"
)

func (a *Activities) CompareCandidates(ctx context.Context, input CandidateCompareInput) (*CandidateCompareOutput, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Comparing candidates started", zap.Int("candidates count", len(input.Candidates)), zap.String("job requirements", input.JobRequirements))

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

	result, err := a.comparator.Compare(ctx, candidates, input.JobRequirements, input.Locale)
	if err != nil {
		logger.Error("Failed to compare candidates", zap.Error(err))
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

	logger.Info("Comparing candidates completed", zap.Int("candidates count", len(input.Candidates)))

	return &output, nil
}
