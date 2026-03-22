package activity

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
)

func (a *Activities) ParseJob(ctx context.Context, input JobParseInput) (*JobParseOutput, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Parsing job started", "locale", input.Locale)

	result, err := a.jobParser.Parse(ctx, input.RawText, input.Locale, input.TeamID)
	if err != nil {
		logger.Error("Failed to parse job", "error", err)
		return nil, fmt.Errorf("failed to parse job: %w", err)
	}

	logger.Info("Parsing job completed", 
		"title", result.Title, 
		"requirements_count", len(result.Requirements), 
		"work_format", result.WorkFormat,
	)

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
