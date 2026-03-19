package hiring

import (
	"context"
	"errors"

	"backend/internal/domain"
	pb "backend/internal/proto/hiring/v1"
	"backend/internal/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetCandidate(ctx context.Context, req *pb.GetCandidateRequest) (*pb.GetCandidateResponse, error) {
	if req.GetCandidateId() == "" {
		return nil, status.Error(codes.InvalidArgument, "candidate_id is required")
	}

	candidate, _, parsedText, _, _, err := h.candidateUC.GetCandidateByID(ctx, req.GetCandidateId())
	if err != nil {
		if errors.Is(err, usecase.ErrCandidateNotFound) {
			return nil, status.Error(codes.NotFound, "candidate not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get candidate: %v", err)
	}

	text := ""
	if parsedText != nil {
		text = *parsedText
	}

	fname := ""
	if candidate.FirstName != nil {
		fname = *candidate.FirstName
	}
	lname := ""
	if candidate.LastName != nil {
		lname = *candidate.LastName
	}

	return &pb.GetCandidateResponse{
		CandidateId: candidate.ID,
		ParsedText:  text,
		FirstName:   fname,
		LastName:    lname,
	}, nil
}

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

