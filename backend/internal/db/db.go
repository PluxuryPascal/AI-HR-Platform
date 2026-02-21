package db

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Client struct {
	log  *zap.Logger
	cfg  *pgxpool.Config
	Pool *pgxpool.Pool

	afterRunFuncs []AfterRun
}

func (c *Client) DependsOn() []string {
	return []string{"logger"}
}

func (c *Client) HealthCheck(ctx context.Context) error {
	if c.Pool == nil {
		return errors.New("database pool is not initialized")
	}

	if err := c.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	c.log.Info("database connection is healthy")

	return nil
}

func (c *Client) Init(ctx context.Context) error {
	pool, err := pgxpool.NewWithConfig(ctx, c.cfg)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}

	for _, f := range c.afterRunFuncs {
		if err := f(ctx, pool); err != nil {
			return fmt.Errorf("failed to run after run function: %w", err)
		}
	}

	c.afterRunFuncs = nil
	c.Pool = pool

	c.log.Info("database initialized")

	return nil
}

func (c *Client) Name() string {
	return "db"
}
func (c *Client) Run(ctx context.Context) error {
	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	<-ctx.Done()

	if c.Pool == nil {
		return fmt.Errorf("database pool is empty")
	}

	c.Pool.Close()

	c.log.Info("database pool closed successfully")

	return nil
}

type AfterRun func(ctx context.Context, pool *pgxpool.Pool) error

func NewDb(log *zap.Logger, conf *config.Config) (*Client, error) {
	cfg, err := pgxpool.ParseConfig(conf.Database.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	cfg.ConnConfig.ConnectTimeout = conf.Database.ConnectionTimeout

	return &Client{
		log:           log,
		cfg:           cfg,
		afterRunFuncs: make([]AfterRun, 0),
	}, nil
}

func (c *Client) AddAfterRun(f ...AfterRun) {
	c.afterRunFuncs = append(c.afterRunFuncs, f...)
}

var _ svc.Service = (*Client)(nil)
