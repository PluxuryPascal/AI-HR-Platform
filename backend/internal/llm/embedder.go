package llm

import (
	"backend/internal/domain"
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const (
	defaultChunkSize    = 512
	defaultChunkOverlap = 128
)

type Embedder struct {
	provider Provider
}

func NewEmbedder(provider Provider) *Embedder {
	return &Embedder{
		provider: provider,
	}
}

func (e *Embedder) EmbedChunks(ctx context.Context, text string, teamID string) ([]domain.EmbeddingChunk, error) {
	cfg := e.provider.GetGlobalConfig()
	chunkSize := cfg.ChunkSize
	if chunkSize == 0 {
		chunkSize = defaultChunkSize
	}
	chunkOverlap := cfg.ChunkOverlap
	if chunkOverlap == 0 {
		chunkOverlap = defaultChunkOverlap
	}

	chunks := splitIntoChunks(text, chunkSize, chunkOverlap)

	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks in text")
	}

	client, settings, err := e.provider.GetClient(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm client: %w", err)
	}

	model := "text-embedding-3-small"
	if settings.EmbedModel != nil && *settings.EmbedModel != "" {
		model = *settings.EmbedModel
	} else if cfg.EmbedModel != "" {
		model = cfg.EmbedModel
	}

	results := make([]domain.EmbeddingChunk, 0, len(chunks))

	for _, chunk := range chunks {
		embedding, err := embedSingle(ctx, client, model, chunk)
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

func embedSingle(ctx context.Context, client *openai.Client, model, text string) ([]float32, error) {
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(model),
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
