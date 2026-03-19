package activity

import (
	"backend/internal/audit"
	"backend/internal/domain"
	"backend/internal/llm"
	"backend/pkg/grpc"
	"backend/pkg/pdf"
	"backend/pkg/storage"
	"context"

	"go.uber.org/zap"
)

type ResumeParser interface {
	Parse(ctx context.Context, text, jobRequirements, locale string, teamID string) (*domain.ParseResult, error)
}

type Scorer interface {
	Score(ctx context.Context, text, jobRequirements, locale string, teamID string) (*domain.ScoreResult, error)
}

type Embedder interface {
	EmbedChunks(ctx context.Context, text string, teamID string) ([]domain.EmbeddingChunk, error)
}

type EmailGenerator interface {
	Generate(ctx context.Context, input llm.EmailGenerateInput) (*llm.EmailGenerateResult, error)
}

type CandidateComparator interface {
	Compare(ctx context.Context, candidates []llm.CandidateCompareInput, jobRequirements, locale string, teamID string) (llm.CompareResult, error)
}

type JobParser interface {
	Parse(ctx context.Context, rawtext, locale string, teamID string) (*llm.JobParseResult, error)
}

type ChatAssistant interface {
	GenerateResponse(ctx context.Context, teamID string, history []domain.ChatMessage, contextChunks []string, question string, locale string) (string, error)
}

type InterviewGenerator interface {
	Generate(ctx context.Context, teamID string, resumeText string, locale string) ([]llm.InterviewQuestion, error)
}

type CandidateStore interface {
	SaveCandidateScore(ctx context.Context, score *domain.CandidateScore, factors []domain.ScoreFactor) error
	SaveResumeEmbedding(ctx context.Context, embedding *domain.ResumeEmbedding) error
	SearchEmbeddings(ctx context.Context, teamID string, candidateID *string, queryVector []float32, limit int) ([]domain.ResumeEmbedding, error)
}

type JobStore interface {
	GetJobRequirements(ctx context.Context, jobID string) ([]byte, error)
}

type CommunicationStore interface {
	Create(ctx context.Context, c *domain.Communication) error
	GetByCandidateID(ctx context.Context, candidateID string) ([]domain.Communication, error)
}

type ChatStore interface {
	AddMessage(ctx context.Context, message *domain.ChatMessage) error
	GetHistory(ctx context.Context, sessionID string) ([]domain.ChatMessage, error)
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
	auditor      *audit.Logger
	chatAsst     ChatAssistant
	chatDB       ChatStore
	interviewGen InterviewGenerator
}

func NewActivities(log *zap.Logger, pdfExtractor pdf.Extractor, storage storage.FileStorage, parser ResumeParser, scorer Scorer, embedder Embedder, emailGen EmailGenerator, comparator CandidateComparator, jobParser JobParser, candidateDB CandidateStore, jobDB JobStore, commDB CommunicationStore, hiringGRPC *grpc.Client, auditor *audit.Logger, chatAsst ChatAssistant, chatDB ChatStore, interviewGen InterviewGenerator) *Activities {
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
		auditor:      auditor,
		chatAsst:     chatAsst,
		chatDB:       chatDB,
		interviewGen: interviewGen,
	}
}
