package hiring

import (
	pb "backend/internal/proto/hiring/v1"
	"backend/internal/usecase"

	"go.uber.org/zap"
)

// Handler реализует gRPC интерфейс HiringServiceServer, выступая в роли
// сборщика (роутера) для различных зон ответственности в рамках сервиса Hiring.
type Handler struct {
	pb.UnimplementedHiringServiceServer

	logger   *zap.Logger
	accessUC usecase.AccessUseCase
}

func NewHandler(log *zap.Logger, accessUC usecase.AccessUseCase) *Handler {
	return &Handler{
		logger:   log,
		accessUC: accessUC,
	}
}
