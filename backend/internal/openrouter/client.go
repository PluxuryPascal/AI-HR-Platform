package openrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	log        *zap.Logger
}

type ModelPricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
	Image      string `json:"image"`
	Request    string `json:"request"`
}

type Architecture struct {
	Modality     string `json:"modality"`
	Tokenizer    string `json:"tokenizer"`
	InstructType string `json:"instruct_type"`
}

type Model struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Pricing       ModelPricing `json:"pricing"`
	ContextLength int          `json:"context_length"`
	Architecture  Architecture `json:"architecture"`
	IsEmbedding   bool         `json:"is_embedding"`
}

type modelsResponse struct {
	Data []Model `json:"data"`
}

func NewClient(log *zap.Logger, baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		log:        log,
	}
}

func (c *Client) GetAvailableModels(ctx context.Context) ([]Model, error) {
	url := fmt.Sprintf("%s/models", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create openrouter request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute openrouter request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter returned status: %d", resp.StatusCode)
	}

	var parsed modelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode openrouter models: %w", err)
	}

	// Post-process to identify embeddings. OpenRouter sometimes lacks a clear "is_embedding" flag, so we check architecture modality or name.
	for i := range parsed.Data {
		m := &parsed.Data[i]
		if m.Architecture.Modality == "text-to-embedding" || m.Architecture.Modality == "embedding" || 
			len(m.Name) > 0 && (containsIgnoreCase(m.Name, "embed") || containsIgnoreCase(m.ID, "embed")) {
			m.IsEmbedding = true
		}
	}

	// Because OpenRouter's /models endpoint does not return embedding-only models, 
	// we manually append common embedding models available on OpenRouter.
	knownEmbeddings := []Model{
		{
			ID: "openai/text-embedding-3-small", 
			Name: "OpenAI: text-embedding-3-small", 
			IsEmbedding: true,
			ContextLength: 8191,
		},
		{
			ID: "openai/text-embedding-3-large", 
			Name: "OpenAI: text-embedding-3-large", 
			IsEmbedding: true,
			ContextLength: 8191,
		},
		{
			ID: "jinaai/jina-embeddings-v2-base-en", 
			Name: "Jina: embeddings-v2-base-en", 
			IsEmbedding: true,
			ContextLength: 8192,
		},
		{
			ID: "nomic-ai/nomic-embed-text-v1.5", 
			Name: "Nomic: embed-text-v1.5", 
			IsEmbedding: true,
			ContextLength: 8192,
		},
	}
	parsed.Data = append(parsed.Data, knownEmbeddings...)

	return parsed.Data, nil
}

func containsIgnoreCase(s, substr string) bool {
	// Simple lowercase check
	sLow := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		sLow += string(c)
	}
	subLow := ""
	for _, c := range substr {
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		subLow += string(c)
	}
	
	for i := 0; i <= len(sLow)-len(subLow); i++ {
		if sLow[i:i+len(subLow)] == subLow {
			return true
		}
	}
	return false
}
