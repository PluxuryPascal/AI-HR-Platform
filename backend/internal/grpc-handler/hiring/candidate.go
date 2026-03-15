package hiring

import (
	"context"

	"backend/internal/domain"
	pb "backend/internal/proto/hiring/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateCandidateProfile(ctx context.Context, req *pb.UpdateCandidateProfileRequest) (*pb.UpdateCandidateProfileResponse, error) {
	if req.GetCandidateId() == "" {
		return nil, status.Error(codes.InvalidArgument, "candidate_id is required")
	}
	if req.GetJobId() == "" {
		return nil, status.Error(codes.InvalidArgument, "job_id is required")
	}

	h.logger.Debug("received UpdateCandidateProfile request", zap.String("candidate_id", req.GetCandidateId()))

	result := h.parseUpdateRequest(req)

	if err := h.candidateUC.FinalizeAIParsing(ctx, result); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to finalize AI parsing: %v", err)
	}

	h.logger.Debug("AI parsing finalized", zap.String("candidate_id", result.CandidateID))

	return &pb.UpdateCandidateProfileResponse{
		Success: true,
	}, nil
}

func (h *Handler) parseUpdateRequest(req *pb.UpdateCandidateProfileRequest) domain.AIParsingResult {
	result := domain.AIParsingResult{
		CandidateID:    req.GetCandidateId(),
		JobID:          req.GetJobId(),
		ParseStatus:    mapParsingStatus(req.GetParsingStatus()),
		StructuredData: req.GetStructuredData(),
	}

	if v := req.GetParsedText(); v != "" {
		result.ParsedText = &v
	}
	if v := req.GetFirstName(); v != "" {
		result.FirstName = &v
	}
	if v := req.GetLastName(); v != "" {
		result.LastName = &v
	}
	if v := req.GetEmail(); v != "" {
		result.Email = &v
	}
	if v := req.GetLocation(); v != "" {
		result.Location = &v
	}
	if skills := req.GetSkills(); len(skills) > 0 {
		result.Skills = skills
	}

	if result.ParseStatus == domain.ParsingStatusNeedsReview {
		result.MissingFields = collectMissingFields(result)
	}

	return result
}

// mapParsingStatus конвертирует proto ParsingStatus → domain CandidateParsingStatus.
func mapParsingStatus(s pb.ParsingStatus) domain.CandidateParsingStatus {
	switch s {
	case pb.ParsingStatus_PARSING_STATUS_SUCCESS:
		return domain.ParsingStatusCompleted
	case pb.ParsingStatus_PARSING_STATUS_NEEDS_REVIEW:
		return domain.ParsingStatusNeedsReview
	default:
		// FAILED и UNSPECIFIED оба → failed
		return domain.ParsingStatusFailed
	}
}

// collectMissingFields возвращает список имён полей, которые AI не смог извлечь.
// Проверяем только поля, критичные для идентификации кандидата.
func collectMissingFields(r domain.AIParsingResult) []string {
	var missing []string

	if r.FirstName == nil {
		missing = append(missing, "first_name")
	}
	if r.LastName == nil {
		missing = append(missing, "last_name")
	}
	if r.Email == nil {
		missing = append(missing, "email")
	}

	return missing
}
