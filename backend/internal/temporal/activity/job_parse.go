package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.uber.org/zap"
)

func (a *Activities) ParseJob(ctx context.Context, input JobParseInput) (*JobParseOutput, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Parsing job started", zap.String("locale", input.Locale))

	result, err := a.jobParser.Parse(ctx, input.RawText, input.Locale)
	if err != nil {
		logger.Error("Failed to parse job", zap.Error(err))
		return nil, fmt.Errorf("failed to parse job: %w", err)
	}

	logger.Info("Parsing job completed", zap.String("title", result.Title), zap.Int("requirements count", len(result.Requirements)), zap.String("work format", result.WorkFormat))

	return &JobParseOutput{
		Title:        result.Title,
		Description:  result.Description,
		Requirements: result.Requirements,
		WorkFormat:   result.WorkFormat,
		SalaryMin:    result.SalaryMin,
		SalaryMax:    result.SalaryMax,
		Currency:     result.Currency,
	}, nil
}
