package temporal

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"backend/pkg/zapadapter"
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type Client struct {
	log            *zap.Logger
	temporalConf   *client.Options
	conf           *config.Temporal
	TemporalClient client.Client
}

func (c *Client) DependsOn() []string {
	return []string{"logger"}
}

func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.TemporalClient.CheckHealth(ctx, &client.CheckHealthRequest{})
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

	c.TemporalClient = cl

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
	if c.TemporalClient != nil {
		c.TemporalClient.Close()
	}

	return nil
}

func (c *Client) StartWorkflow(
	ctx context.Context,
	workflowID string,
	timeout time.Duration,
	workflowFn any,
	input any,
) (string, error) {
	opts := client.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                c.conf.QueueName,
		WorkflowExecutionTimeout: timeout,
		WorkflowIDConflictPolicy: 1,
	}

	we, err := c.TemporalClient.ExecuteWorkflow(ctx, opts, workflowFn, input)
	if err != nil {
		return "", fmt.Errorf("start workflow %q: %w", workflowID, err)
	}

	c.log.Info("workflow started",
		zap.String("workflow_id", we.GetID()),
		zap.String("run_id", we.GetRunID()),
	)

	return we.GetID(), nil
}

func (c *Client) CancelWorkflow(ctx context.Context, workflowID string) error {
	if err := c.TemporalClient.CancelWorkflow(ctx, workflowID, ""); err != nil {
		return fmt.Errorf("cancel workflow %s: %w", workflowID, err)
	}

	c.log.Info("workflow cancelled", zap.String("workflow_id", workflowID))

	return nil
}

func (c *Client) TaskQueue() string {
	return c.conf.QueueName
}

var _ svc.Service = (*Client)(nil)

func NewClient(log *zap.Logger, temporalConf *client.Options, conf *config.Temporal) *Client {
	return &Client{
		log: log,
		temporalConf: &client.Options{
			HostPort:  conf.HostPort,
			Namespace: conf.Namespace,
			Logger:    zapadapter.NewZapAdapter(log),
		},
		conf: conf,
	}
}
