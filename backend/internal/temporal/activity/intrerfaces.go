package activity

import (
	"backend/internal/domain"
	"backend/internal/llm"
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

type EmailGenerator interface {
	Generate(ctx context.Context, input llm.EmailGenerateInput) (*llm.EmailGenerateResult, error)
}

type CandidateComparator interface {
	Compare(ctx context.Context, candidates []llm.CandidateCompareInput, jobRequirements string) (llm.CompareResult, error)
}

type JobParser interface {
	Parse(ctx context.Context, rawtext, locale string) (*llm.JobParseResult, error)
}

type CandidateStore interface {
	SaveCandidateScore(ctx context.Context, score *domain.CandidateScore, factors []domain.ScoreFactor) error
	SaveResumeEmbedding(ctx context.Context, embedding *domain.ResumeEmbedding) error
}

type JobStore interface {
	GetJobRequirements(ctx context.Context, jobID string) ([]byte, error)
}

type CommunicationStore interface {
	Create(ctx context.Context, c *domain.Communication) error
	GetByCandidateID(ctx context.Context, candidateID string) ([]domain.Communication, error)
}

type Activities struct {
	log          *zap.Logger
	pdfExtractor pdf.Extractor
	storage      storage.FileStorage
	parser       ResumeParser
	scorer       Scorer
	embedder     Embedder
	emailGen     EmailGenerator
	comparator   CandidateComparator
	jobParser    JobParser
	candidateDB  CandidateStore
	jobDB        JobStore
	commDB       CommunicationStore
	hiringGRPC   *grpc.Client
}

func NewActivities(log *zap.Logger, pdfExtractor pdf.Extractor, storage storage.FileStorage, parser ResumeParser, scorer Scorer, embedder Embedder, emailGen EmailGenerator, comparator CandidateComparator, jobParser JobParser, candidateDB CandidateStore, jobDB JobStore, commDB CommunicationStore, hiringGRPC *grpc.Client) *Activities {
	return &Activities{
		log:          log,
		pdfExtractor: pdfExtractor,
		storage:      storage,
		parser:       parser,
		scorer:       scorer,
		embedder:     embedder,
		emailGen:     emailGen,
		comparator:   comparator,
		jobParser:    jobParser,
		candidateDB:  candidateDB,
		jobDB:        jobDB,
		commDB:       commDB,
		hiringGRPC:   hiringGRPC,
	}
}
