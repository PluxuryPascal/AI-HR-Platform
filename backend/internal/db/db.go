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

type PostgresClient struct {
	log  *zap.Logger
	cfg  *pgxpool.Config
	Pool *pgxpool.Pool

	afterRunFuncs []AfterRun
}

func (c *PostgresClient) DependsOn() []string {
	return []string{"logger"}
}

func (c *PostgresClient) HealthCheck(ctx context.Context) error {
	if c.Pool == nil {
		return errors.New("database pool is not initialized")
	}

	if err := c.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	c.log.Debug("database connection is healthy")

	return nil
}

func (c *PostgresClient) Init(ctx context.Context) error {
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

	c.log.Debug("database initialized")

	return nil
}

func (c *PostgresClient) Name() string {
	return "db"
}
func (c *PostgresClient) Run(ctx context.Context) error {
	return nil
}

func (c *PostgresClient) Stop(ctx context.Context) error {
	if c.Pool == nil {
		return fmt.Errorf("database pool is empty")
	}

	c.Pool.Close()

	c.log.Debug("database pool closed successfully")

	return nil
}

type AfterRun func(ctx context.Context, pool *pgxpool.Pool) error

func NewDb(log *zap.Logger, conf *config.Config) (*PostgresClient, error) {
	cfg, err := pgxpool.ParseConfig(conf.Database.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	cfg.ConnConfig.ConnectTimeout = conf.Database.ConnectionTimeout
	cfg.MaxConns = conf.Database.MaxConns
	cfg.MinConns = conf.Database.MinConns
	cfg.MaxConnLifetime = conf.Database.MaxConnLifetime

	return &PostgresClient{
		log:           log,
		cfg:           cfg,
		afterRunFuncs: make([]AfterRun, 0),
	}, nil
}

func (c *PostgresClient) AddAfterRun(f ...AfterRun) {
	c.afterRunFuncs = append(c.afterRunFuncs, f...)
}

var _ svc.Service = (*PostgresClient)(nil)
