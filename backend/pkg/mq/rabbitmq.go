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

// RabbitMQ реализует svc.Service и управляет единственным долгоживущим
// AMQP-соединением для всего приложения.
//
// # Почему одно соединение на всё?
//
// RabbitMQ рекомендует создавать одно TCP-соединение (*Conn) и
// мультиплексировать его между Publisher/Consumer через виртуальные
// каналы (Channel). Создание отдельных TCP-соединений для каждого
// Publisher-а — расточительно и не масштабируется.
//
// go-rabbitmq автоматически управляет переподключением при обрыве сети:
// и Publisher, и Consumer восстановятся без ручного вмешательства.
//
// # Ответственность
//
//   - Установить соединение и передать *Conn в Publisher и Consumer.
//   - Корректно закрыть соединение при остановке (Stop).
//
// # Что НЕ делает этот сервис
//
// Объявление очередей и Exchange-ов — НЕ здесь. Каждый Consumer
// объявляет нужную ему очередь самостоятельно при старте (через
// WithConsumerOptionsQueueArgs и другие опции). Это правильный подход:
// Consumer «владеет» своей очередью и знает её параметры.
type RabbitMQ struct {
	log *zap.Logger
	cfg *config.RabbitMQ

	// Conn — публичное поле: передаётся в NewMQPublisher и NewMQConsumer.
	Conn *rabbitmq.Conn
}

var _ svc.Service = (*RabbitMQ)(nil)

func NewRabbitMQ(log *zap.Logger, cfg *config.RabbitMQ) *RabbitMQ {
	return &RabbitMQ{log: log, cfg: cfg}
}

func (r *RabbitMQ) Name() string        { return "rabbitmq" }
func (r *RabbitMQ) DependsOn() []string { return []string{"logger"} }

// Init устанавливает TCP-соединение с RabbitMQ-брокером.
//
// go-rabbitmq принимает строку URL вида "amqp://user:pass@host:5672/" и
// самостоятельно запускает фоновую горутину для автоматического
// переподключения с интервалом cfg.ReconnectDelay.
func (r *RabbitMQ) Init(_ context.Context) error {
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

// HealthCheck убеждается, что Init был вызван успешно.
func (r *RabbitMQ) HealthCheck(_ context.Context) error {
	if r.Conn == nil {
		return errors.New("rabbitmq connection is not initialized")
	}

	return nil
}

// Run — no-op: go-rabbitmq управляет переподключением в своей горутине.
func (r *RabbitMQ) Run(_ context.Context) error { return nil }

// Stop закрывает TCP-соединение. Перед этим svc.Run уже остановил
// все Consumer и Publisher (они объявляют зависимость DependsOn: "rabbitmq",
// поэтому останавливаются первыми).
func (r *RabbitMQ) Stop(_ context.Context) error {
	if err := r.Conn.Close(); err != nil {
		return fmt.Errorf("failed to close rabbitmq connection: %w", err)
	}

	return nil
}
