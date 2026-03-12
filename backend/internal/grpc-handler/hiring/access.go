package hiring

import (
	"context"

	pb "backend/internal/proto/hiring/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) TransferJobAccess(ctx context.Context, req *pb.TransferJobAccessRequest) (*pb.TransferJobAccessResponse, error) {
	userID := req.GetUserId()
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	jobIDs := req.GetJobIds()
	if len(jobIDs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "job_ids cannot be empty")
	}

	h.logger.Debug("received TransferJobAccess request",
		zap.String("user_id", userID),
		zap.Int("jobs_count", len(jobIDs)),
	)

	err := h.accessUC.TransferAccess(ctx, userID, jobIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transfer job access: %v", err)
	}

	return &pb.TransferJobAccessResponse{
		Success: true,
	}, nil
}
