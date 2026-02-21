package db

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	log  *zap.Logger
	cfg  *redis.Options
	Pool *redis.Client
}

func (r *RedisClient) DependsOn() []string {
	return []string{"logger"}
}

func (r *RedisClient) HealthCheck(ctx context.Context) error {
	if r.Pool == nil {
		return fmt.Errorf("redis pool is not initialized")
	}

	if err := r.Pool.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	r.log.Info("redis connection is healthy")

	return nil
}

func (r *RedisClient) Init(ctx context.Context) error {
	pool := redis.NewClient(r.cfg)
	if pool == nil {
		return fmt.Errorf("in memory pool is nil")
	}

	r.Pool = pool

	r.log.Info("redis initialized")

	return nil
}

func (r *RedisClient) Name() string {
	return "redis"
}

func (r *RedisClient) Run(ctx context.Context) error {
	return nil
}

func (r *RedisClient) Stop(ctx context.Context) error {
	<-ctx.Done()

	if r.Pool == nil {
		return fmt.Errorf("redis pool is empty")
	}

	if err := r.Pool.Close(); err != nil {
		return fmt.Errorf("close redis pool: %w", err)
	}

	r.log.Info("redis pool closed successfully")

	return nil
}

func NewRedis(log *zap.Logger, conf *config.Config) (*RedisClient, error) {
	cfg, err := redis.ParseURL(conf.Redis.Host)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	cfg.DialTimeout = conf.Redis.ConnectionTimeout

	return &RedisClient{
		log: log,
		cfg: cfg,
	}, nil
}

var _ svc.Service = (*RedisClient)(nil)
