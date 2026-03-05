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

type MQPublisher struct {
	log *zap.Logger
	cfg *config.RabbitMQ

	conn *rabbitmq.Conn
	Pub  *rabbitmq.Publisher
}

func (m *MQPublisher) DependsOn() []string {
	return []string{"logger", "rabbitmq"}
}

func (m *MQPublisher) HealthCheck(ctx context.Context) error {
	if m.Pub == nil {
		return errors.New("publisher is not initialized")
	}

	return nil
}

func (m *MQPublisher) Init(ctx context.Context) error {
	pub, err := rabbitmq.NewPublisher(
		m.conn,
		rabbitmq.WithPublisherOptionsLogger(newZapLogger(m.log)),
	)
	if err != nil {
		return fmt.Errorf("failed to create rabbitmq publisher: %w", err)
	}

	m.Pub = pub

	m.log.Debug("rabbitmq publisher initialized")

	return nil
}

func (m *MQPublisher) Name() string {
	return "rabbitmq_publisher"
}

func (m *MQPublisher) Run(ctx context.Context) error {
	return nil
}
func (m *MQPublisher) Stop(ctx context.Context) error {
	m.Pub.Close()

	return nil
}

var _ svc.Service = (*MQPublisher)(nil)

func NewMQPublisher(log *zap.Logger, cfg *config.RabbitMQ, conn *rabbitmq.Conn) *MQPublisher {
	return &MQPublisher{
		log:  log,
		cfg:  cfg,
		conn: conn,
	}
}
