package activity

import (
	"backend/internal/audit"
	"backend/internal/domain"
	hiringv1 "backend/internal/proto/hiring/v1"
	"backend/pkg/pdf"
	"context"
	"errors"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/zap"
)

func (a *Activities) PDFExtract(ctx context.Context, input ResumePipelineInput) (*string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("PDF extraction started", zap.String("candidate id", input.CandidateID), zap.String("file key", input.ResumeFileKey))

	file, err := a.storage.DownloadFile(ctx, input.ResumeFileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	parsedText, err := a.pdfExtractor.Extract(file)
	if err != nil {
		switch {
		case errors.Is(err, pdf.ErrPasswordProtected),
			errors.Is(err, pdf.ErrEmptyContent),
			errors.Is(err, pdf.ErrImageBasedPDF):
			return nil, temporal.NewApplicationError(
				err.Error(),
				"NonRetryablePDFError",
				err,
			)
		}
		return nil, fmt.Errorf("extract pdf text: %w", err)
	}

	logger.Info("PDF extraction completed", zap.String("candidate id", input.CandidateID))

	return &parsedText, nil
}

func (a *Activities) LLMParse(ctx context.Context, input LLMParseInput) (*domain.ParseResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("LLM parsing started", zap.String("candidate id", input.CandidateID), zap.String("job id", input.JobID))

	var jobRequirements string
	jobRequirementsBytes, err := a.jobDB.GetJobRequirements(ctx, input.JobID)
	if err != nil {
		logger.Warn("Failed to get job requirements", zap.Error(err))
	} else {
		jobRequirements = string(jobRequirementsBytes)
	}

	result, err := a.parser.Parse(ctx, input.ParsedText, jobRequirements, input.Locale)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	result.JobRequirements = jobRequirements

	logger.Info("LLM parsing completed", zap.String("candidate id", input.CandidateID), zap.String("status", string(result.ParsingStatus())))

	return result, nil
}

func (a *Activities) LLMScore(ctx context.Context, input LLMScoreInput) (*domain.ScoreResult, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("LLMScore started", "candidate_id", input.CandidateID)

	result, err := a.scorer.Score(ctx, input.ParsedText, input.JobRequirements, input.Locale)
	if err != nil {
		return nil, fmt.Errorf("failed to score: %w", err)
	}

	score := &domain.CandidateScore{
		CandidateID: input.CandidateID,
		MatchScore:  result.MatchScore,
	}

	factors := make([]domain.ScoreFactor, len(result.Factors))
	for i, factor := range result.Factors {
		factors[i] = domain.ScoreFactor{
			CandidateID: input.CandidateID,
			Type:        factor.Type,
			Description: factor.Description,
			Impact:      factor.Impact,
		}
	}

	if err := a.candidateDB.SaveCandidateScore(ctx, score, factors); err != nil {
		return nil, fmt.Errorf("failed to save score: %w", err)
	}

	_ = a.auditor.Log(ctx, audit.Entry{
		TeamID:    input.TeamID,
		ActorType: audit.ActorSystem,
		Action:    audit.AIScoreComputed,
		TargetID:  &input.CandidateID,
	})

	logger.Info("LLMScore completed", zap.String("candidate id", input.CandidateID), zap.Int("match score", result.MatchScore))

	return result, nil
}

func (a *Activities) LLMEmbed(ctx context.Context, input LLMEmbedInput) error {
	logger := activity.GetLogger(ctx)

	logger.Info("LLMEmbed started", "candidate_id", input.CandidateID)

	chunks, err := a.embedder.EmbedChunks(ctx, input.ParsedText)
	if err != nil {
		return fmt.Errorf("failed to embed: %w", err)
	}

	for _, chunk := range chunks {
		emb := &domain.ResumeEmbedding{
			CandidateID: input.CandidateID,
			TeamID:      input.TeamID,
			Embedding:   chunk.Embedding,
			ChunkText:   chunk.Text,
		}

		if err := a.candidateDB.SaveResumeEmbedding(ctx, emb); err != nil {
			return fmt.Errorf("failed to save embedding: %w", err)
		}
	}

	logger.Info("LLMEmbed completed", zap.String("candidate id", input.CandidateID))

	return nil
}

func (a *Activities) GRPCCallback(ctx context.Context, input GRPCCallbackInput, teamID string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("GRPCCallback started",
		"candidate_id", input.CandidateID,
		"status", input.Status.String(),
	)

	hiringClient := hiringv1.NewHiringServiceClient(a.hiringGRPC.GetConn())

	req := &hiringv1.UpdateCandidateProfileRequest{
		CandidateId:   input.CandidateID,
		JobId:         input.JobID,
		ParsingStatus: input.Status,
	}

	if input.Status != hiringv1.ParsingStatus_PARSING_STATUS_FAILED {
		pr := input.ParseResult
		req.ParsedText = input.ParsedText
		req.FirstName = ptrVal(pr.FirstName)
		req.LastName = ptrVal(pr.LastName)
		req.Email = ptrVal(pr.Email)
		req.Location = ptrVal(pr.Location)
		req.Skills = pr.Skills
		req.StructuredData = pr.StructuredData
	}

	if _, err := hiringClient.UpdateCandidateProfile(ctx, req); err != nil {
		return fmt.Errorf("update candidate profile rpc: %w", err)
	}

	if input.Status != hiringv1.ParsingStatus_PARSING_STATUS_FAILED {
		_ = a.auditor.Log(ctx, audit.Entry{
			TeamID:    teamID,
			ActorType: audit.ActorSystem,
			Action:    audit.AIResumeParsed,
			TargetID:  &input.CandidateID,
		})
	}

	logger.Info("GRPCCallback completed",
		"candidate_id", input.CandidateID,
		"status", input.Status.String(),
	)

	return nil
}

func ptrVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
