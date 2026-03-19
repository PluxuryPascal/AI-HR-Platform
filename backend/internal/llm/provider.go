package llm

import (
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/pkg/config"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type Provider interface {
	// GetClient returns an initialized openai.Client and the team's settings.
	// Returns an error if the team has no API key configured.
	GetClient(ctx context.Context, teamID string) (*openai.Client, *domain.TeamAISettings, error)
	// GetGlobalConfig returns the global openrouter configuration (timeouts, limits, etc.)
	GetGlobalConfig() *config.OpenRouter
}

type provider struct {
	log         *zap.Logger
	cfg         *config.OpenRouter
	settingsRepo repo.AiSettingsRepository
}

func NewProvider(log *zap.Logger, cfg *config.OpenRouter, settingsRepo repo.AiSettingsRepository) Provider {
	return &provider{
		log:          log,
		cfg:          cfg,
		settingsRepo: settingsRepo,
	}
}

func (p *provider) GetGlobalConfig() *config.OpenRouter {
	return p.cfg
}

func (p *provider) GetClient(ctx context.Context, teamID string) (*openai.Client, *domain.TeamAISettings, error) {
	if teamID == "" {
		return nil, nil, fmt.Errorf("llm provider: team_id is required")
	}

	settings, err := p.settingsRepo.GetByTeamID(ctx, teamID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ai settings: %w", err)
	}

	if settings == nil || settings.APIKey == nil || *settings.APIKey == "" {
		return nil, nil, fmt.Errorf("ai settings or api key not configured for team: %s", teamID)
	}

	timeout := p.cfg.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	clientCfg := openai.DefaultConfig(*settings.APIKey)
	clientCfg.BaseURL = p.cfg.BaseURL
	clientCfg.HTTPClient = &http.Client{
		Timeout: timeout,
	}

	client := openai.NewClientWithConfig(clientCfg)

	return client, settings, nil
}
