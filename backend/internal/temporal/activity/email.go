package activity

import (
	"backend/internal/domain"
	"backend/internal/llm"
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.uber.org/zap"
)

func (a *Activities) GenerateEmail(ctx context.Context, input EmailGenerateInput) (*EmailGenerateOutput, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Generating email started", zap.String("candidate id", input.CandidateID), zap.String("locale", input.Locale))

	result, err := a.emailGen.Generate(ctx, llm.EmailGenerateInput{
		CandidateFirstName: input.CandidateFirstName,
		CandidateLastName:  input.CandidateLastName,
		CandidateEmail:     input.CandidateEmail,
		Role:               input.Role,
		Skills:             input.Skills,
		MatchScore:         input.MatchScore,
		Type:               domain.EmailType(input.EmailType),
		Tone:               llm.EmailTone(input.Tone),
		Locale:             input.Locale,
		RecruiterName:      input.RecruiterName,
		CompanyName:        input.CompanyName,
	})
	if err != nil {
		logger.Error("Failed to generate email", zap.Error(err))
		return nil, fmt.Errorf("failed to generate email: %w", err)
	}

	logger.Info("Generating email completed", zap.String("candidate id", input.CandidateID))

	return &EmailGenerateOutput{
		Subject: result.Subject,
		Body:    result.Body,
	}, nil
}
