package hiring

import (
	"context"
	"time"

	pb "backend/internal/proto/hiring/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateCandidateProfile(ctx context.Context, req *pb.UpdateCandidateProfileRequest) (*pb.UpdateCandidateProfileResponse, error) {
	candidateID := req.GetCandidateId()
	if candidateID == "" {
		return nil, status.Error(codes.InvalidArgument, "candidate_id is required")
	}

	h.logger.Debug("received UpdateCandidateProfile request", zap.String("candidate_id", candidateID))

	// 1. Get existing candidate and profile
	cand, profile, _, err := h.candidateUC.GetCandidateByID(ctx, candidateID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "candidate not found: %v", err)
	}

	// 2. Update Candidate Basic Info
	if req.FirstName != "" {
		cand.FirstName = &req.FirstName
	}
	if req.LastName != "" {
		cand.LastName = &req.LastName
	}
	if req.Email != "" {
		cand.Email = &req.Email
	}
	if req.ParsedText != "" {
		cand.ParsedText = &req.ParsedText
	}

	if err := h.candidateUC.UpdateCandidate(ctx, cand); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update basic candidate info: %v", err)
	}

	// 3. Update Candidate Profile (Structured Data)
	profile.StructuredData = req.StructuredData
	now := time.Now()
	profile.AIParsedAt = &now

	if err := h.candidateUC.UpdateCandidateProfile(ctx, profile); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update structured profile data: %v", err)
	}

	return &pb.UpdateCandidateProfileResponse{
		Success: true,
	}, nil
}
