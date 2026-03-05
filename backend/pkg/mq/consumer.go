package mq

import (
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"errors"
	"fmt"

	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type MQConsumer struct {
	log *zap.Logger
	cfg *config.RabbitMQ

	conn     *rabbitmq.Conn
	Consumer *rabbitmq.Consumer
}

func (m *MQConsumer) DependsOn() []string {
	return []string{"logger", "rabbitmq"}
}

func (m *MQConsumer) HealthCheck(ctx context.Context) error {
	if m.Consumer == nil {
		return errors.New("consumer is not initialized")
	}

	return nil
}

func (m *MQConsumer) Init(ctx context.Context) error {
	consumer, err := rabbitmq.NewConsumer(
		m.conn,
		m.cfg.Queue,
	)
	if err != nil {
		return fmt.Errorf("failed to create rabbitmq consumer: %w", err)
	}

	m.Consumer = consumer

	m.Consumer.Run()

	m.log.Debug("rabbitmq consumer initialized")

	return nil
}

func (m *MQConsumer) Name() string {
	return "rabbitmq_consumer"
}

func (m *MQConsumer) Run(ctx context.Context) error {
	return nil
}
func (m *MQConsumer) Stop(ctx context.Context) error {
	m.Consumer.Close()

	return nil
}

var _ svc.Service = (*MQConsumer)(nil)

func NewMQConsumer(log *zap.Logger, cfg *config.RabbitMQ, conn *rabbitmq.Conn) *MQConsumer {
	return &MQConsumer{
		log:  log,
		cfg:  cfg,
		conn: conn,
	}
}
