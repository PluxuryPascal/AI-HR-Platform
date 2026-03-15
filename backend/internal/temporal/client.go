package temporal

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"backend/pkg/zapadapter"
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type Client struct {
	log           *zap.Logger
	temporalConf  *client.Options
	conf          *config.Temporal
	TemporaClient client.Client
}

func (c *Client) DependsOn() []string {
	return []string{"logger"}
}

func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.TemporaClient.CheckHealth(ctx, &client.CheckHealthRequest{})
	if err != nil {
		return fmt.Errorf("temporal client is not healthy: %w", err)
	}

	c.log.Info("temporal client is healthy")

	return nil
}

func (c *Client) Init(ctx context.Context) error {
	cl, err := client.Dial(*c.temporalConf)
	if err != nil {
		return err
	}

	c.TemporaClient = cl

	c.log.Info("temporal client initialized")

	return nil
}

func (c *Client) Name() string {
	return "temporal-client"
}

func (c *Client) Run(ctx context.Context) error {
	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	if c.TemporaClient != nil {
		c.TemporaClient.Close()
	}

	return nil
}

var _ svc.Service = (*Client)(nil)

func NewClient(log *zap.Logger, temporalConf *client.Options, conf *config.Temporal) *Client {
	temporalLogger := zapadapter.NewZapAdapter(log)

	return &Client{
		log: log,
		temporalConf: &client.Options{
			HostPort: conf.HostPort,
			Logger:   temporalLogger,
		},
		conf: conf,
	}
}
