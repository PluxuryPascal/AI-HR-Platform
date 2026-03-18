package ai_engine

import (
	"backend/internal/domain"
	pb "backend/internal/proto/ai_engine/v1"
	"backend/internal/temporal"
	"context"

	"go.uber.org/zap"
)

type CandidateRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Candidate, *domain.CandidateProfile, *string, error)
	GetScoreByCandidateID(ctx context.Context, candidateID string) (*domain.CandidateScore, []domain.ScoreFactor, error)
	GetScoresByCandidateIDs(ctx context.Context, candidateIDs []string) (map[string]*domain.CandidateScore, error)
}

type CommunicationRepository interface {
	Create(ctx context.Context, c *domain.Communication) error
	GetByCandidateID(ctx context.Context, candFinalizeAIParsingidateID string) ([]domain.Communication, error)
}

type Handler struct {
	pb.UnimplementedAIEngineServiceServer

	logger         *zap.Logger
	temporalClient *temporal.Client
	candidateRepo  CandidateRepository
	commRepo       CommunicationRepository
}

func NewHandler(
	logger *zap.Logger,
	temporalClient *temporal.Client,
	candidateRepo CandidateRepository,
	commRepo CommunicationRepository,
) *Handler {
	return &Handler{
		logger:         logger,
		temporalClient: temporalClient,
		candidateRepo:  candidateRepo,
		commRepo:       commRepo,
	}
}
