package llm

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type Client struct {
	log    *zap.Logger
	cfg    *config.OpenRouter
	OpenAI *openai.Client
}

func (c *Client) DependsOn() []string {
	return []string{"logger"}
}

func (c *Client) HealthCheck(ctx context.Context) error {
	if c.OpenAI == nil {
		return fmt.Errorf("llm client is not initialized")
	}

	_, err := c.OpenAI.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("llm client health check error: %w", err)
	}

	c.log.Info("llm client health check success")

	return nil
}

func (c *Client) Init(ctx context.Context) error {
	timeout := c.cfg.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	clientCfg := openai.DefaultConfig(c.cfg.APIKey)
	clientCfg.BaseURL = c.cfg.BaseURL
	clientCfg.HTTPClient = &http.Client{
		Timeout: timeout,
	}

	c.OpenAI = openai.NewClientWithConfig(clientCfg)

	c.log.Info("llm client initialized", zap.String("base_url", c.cfg.BaseURL))

	return nil
}

func (c *Client) Name() string {
	return "llm"
}

func (c *Client) Run(ctx context.Context) error {
	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	return nil
}

var _ svc.Service = (*Client)(nil)

func NewClient(log *zap.Logger, cfg *config.OpenRouter) *Client {
	return &Client{
		log: log,
		cfg: cfg,
	}
}
