package hiring

import (
	"context"

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

	// Вызов usecase
	// err := h.candidateUC.UpdateProfile(...)
	// if err != nil {
	//     return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	// }

	return &pb.UpdateCandidateProfileResponse{
		Success: true,
	}, nil
}
