package usecase

import (
	"backend/internal/domain"
	pbAI "backend/internal/proto/ai_engine/v1"
	"context"
	"fmt"
)

type AiUseCase interface {
	ParseJobDescription(ctx context.Context, rawText, locale string, teamID string) (*domain.JobParseResult, error)
	GenerateCandidateEmail(ctx context.Context, input *pbAI.GenerateCandidateEmailRequest) (*pbAI.GenerateCandidateEmailResponse, error)
}

type aiUseCase struct {
	aiClient *pbAI.AIEngineServiceClient
}

func NewAiUseCase(aiClient *pbAI.AIEngineServiceClient) AiUseCase {
	return &aiUseCase{
		aiClient: aiClient,
	}
}

func (u *aiUseCase) ParseJobDescription(ctx context.Context, rawText, locale string, teamID string) (*domain.JobParseResult, error) {
	if u.aiClient == nil || *u.aiClient == nil {
		return nil, fmt.Errorf("ai client not initialized")
	}

	resp, err := (*u.aiClient).ParseJobDescription(ctx, &pbAI.ParseJobDescriptionRequest{
		RawText: rawText,
		Locale:  locale,
		TeamId:  teamID,
	})
	if err != nil {
		return nil, fmt.Errorf("ai engine parse job: %w", err)
	}

	return &domain.JobParseResult{
		Title:        resp.GetTitle(),
		Description:  resp.GetDescription(),
		Requirements: resp.GetRequirements(),
		WorkFormat:   resp.GetWorkFormat(),
		SalaryMin:    int(resp.GetSalaryMin()),
		SalaryMax:    int(resp.GetSalaryMax()),
		Currency:     resp.GetCurrency(),
	}, nil
}

func (u *aiUseCase) GenerateCandidateEmail(ctx context.Context, input *pbAI.GenerateCandidateEmailRequest) (*pbAI.GenerateCandidateEmailResponse, error) {
	if u.aiClient == nil || *u.aiClient == nil {
		return nil, fmt.Errorf("ai client not initialized")
	}

	return (*u.aiClient).GenerateCandidateEmail(ctx, input)
}
