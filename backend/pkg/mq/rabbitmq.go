package mq

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"

	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	log *zap.Logger
	cfg *config.RabbitMQ

	Conn *rabbitmq.Conn
}

func (r *RabbitMQ) DependsOn() []string {
	return []string{"logger"}
}

func (r *RabbitMQ) HealthCheck(ctx context.Context) error {
	return nil
}

func (r *RabbitMQ) Init(ctx context.Context) error {
	conn, err := rabbitmq.NewConn(
		r.cfg.URL,
		rabbitmq.WithConnectionOptionsLogger(newZapLogger(r.log)),
		rabbitmq.WithConnectionOptionsReconnectInterval(r.cfg.ReconnectDelay),
	)
	if err != nil {
		return fmt.Errorf("failed to create rabbitmq connection: %w", err)
	}

	r.Conn = conn

	r.log.Debug("rabbitmq connection established")

	return nil
}

func (r *RabbitMQ) Name() string {
	return "rabbitmq"
}

func (r *RabbitMQ) Run(ctx context.Context) error {
	return nil
}

func (r *RabbitMQ) Stop(ctx context.Context) error {
	if err := r.Conn.Close(); err != nil {
		return fmt.Errorf("failed to close rabbitmq connection: %w", err)
	}

	return nil
}

var _ svc.Service = (*RabbitMQ)(nil)

func NewRabbitMQ(log *zap.Logger, cfg *config.RabbitMQ) *RabbitMQ {
	return &RabbitMQ{
		log: log,
		cfg: cfg,
	}
}
