package activity

import (
	"backend/internal/domain"
	"backend/pkg/grpc"
	"backend/pkg/pdf"
	"backend/pkg/storage"
	"context"

	"go.uber.org/zap"
)

type ResumeParser interface {
	Parse(ctx context.Context, text, jobRequirements string) (*domain.ParseResult, error)
}

type Scorer interface {
	Score(ctx context.Context, text, jobRequirements string) (*domain.ScoreResult, error)
}

type Embedder interface {
	EmbedChunks(ctx context.Context, text string) ([]domain.EmbeddingChunk, error)
}

type CandidateStore interface {
	SaveCandidateScore(ctx context.Context, score *domain.CandidateScore, factors []domain.ScoreFactor) error
	SaveResumeEmbedding(ctx context.Context, embedding *domain.ResumeEmbedding) error
}

type JobStore interface {
	GetJobRequirements(ctx context.Context, jobID string) ([]byte, error)
}

type Activities struct {
	log          *zap.Logger
	pdfExtractor pdf.Extractor
	storage      storage.FileStorage
	parser       ResumeParser
	scorer       Scorer
	embedder     Embedder
	candidateDB  CandidateStore
	jobDB        JobStore
	hiringGRPC   *grpc.Client
}

func NewActivities(log *zap.Logger, pdfExtractor pdf.Extractor, storage storage.FileStorage, parser ResumeParser, scorer Scorer, embedder Embedder, candidateDB CandidateStore, jobDB JobStore, hiringGRPC *grpc.Client) *Activities {
	return &Activities{
		log:          log,
		pdfExtractor: pdfExtractor,
		storage:      storage,
		parser:       parser,
		scorer:       scorer,
		embedder:     embedder,
		candidateDB:  candidateDB,
		jobDB:        jobDB,
		hiringGRPC:   hiringGRPC,
	}
}
