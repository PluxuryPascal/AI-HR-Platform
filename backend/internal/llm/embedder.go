package llm

import (
	"backend/internal/domain"
	"backend/pkg/config"
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const (
	defaultChunkSize    = 512
	defaultChunkOverlap = 128
)

type Embedder struct {
	client       *openai.Client
	cfg          *config.OpenRouter
	chunkSize    int
	chunkOverlap int
}

func NewEmbedder(client *openai.Client, cfg *config.OpenRouter) *Embedder {
	chunkSize := cfg.ChunkSize
	if chunkSize == 0 {
		chunkSize = defaultChunkSize
	}

	chunkOverlap := cfg.ChunkOverlap
	if chunkOverlap == 0 {
		chunkOverlap = defaultChunkOverlap
	}

	return &Embedder{
		client:       client,
		cfg:          cfg,
		chunkSize:    chunkSize,
		chunkOverlap: chunkOverlap,
	}
}

func (e *Embedder) EmbedChunks(ctx context.Context, text string) ([]domain.EmbeddingChunk, error) {
	chunks := splitIntoChunks(text, e.chunkSize, e.chunkOverlap)

	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks in text")
	}

	results := make([]domain.EmbeddingChunk, 0, len(chunks))

	for _, chunk := range chunks {
		embedding, err := e.embedSingle(ctx, chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to create embeddings: %w", err)
		}

		results = append(results, domain.EmbeddingChunk{
			Text:      chunk,
			Embedding: embedding,
		})
	}

	return results, nil
}

func (e *Embedder) embedSingle(ctx context.Context, text string) ([]float32, error) {
	resp, err := e.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(e.cfg.EmbedModel),
		Input: []string{text},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings in response")
	}

	return resp.Data[0].Embedding, nil
}

func splitIntoChunks(text string, chunkSize, overlap int) []string {
	if len(text) == 0 {
		return nil
	}

	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0

	for start < len(text) {
		end := start + chunkSize

		if end >= len(text) {
			chunks = append(chunks, text[start:])
			break
		}

		boundary := findWordBoundary(text, end)
		chunks = append(chunks, text[start:boundary])

		next := boundary - overlap
		if next <= start {
			next = start + 1
		}
		start = next
	}

	return chunks
}

func findWordBoundary(text string, pos int) int {
	if pos >= len(text) {
		return len(text)
	}

	limit := pos - 100
	if limit < 0 {
		limit = 0
	}

	for i := pos; i > limit; i-- {
		if text[i] == ' ' || text[i] == '\n' {
			return i
		}
	}

	return pos
}
