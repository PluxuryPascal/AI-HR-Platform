package mq

import (
	"backend/pkg/svc"
	"context"
	"errors"
	"fmt"

	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

// MQPublisher реализует svc.Service и оборачивает go-rabbitmq Publisher.
//
// Ключевое отличие от первоначальной версии — Publisher не знает заранее,
// в какой exchange и с каким routing key публиковать. Эти параметры
// передаются при каждом вызове Publish() через []PublisherOption.
//
// Это позволяет иметь один экземпляр *MQPublisher и использовать его
// для публикации в РАЗНЫЕ exchange / с разными routing key:
//
//	publisher.Publish(ctx, body,
//	    mq.WithExchange("hiring.events"),
//	    mq.WithRoutingKey("candidate.created"),
//	)
//
//	publisher.Publish(ctx, body,
//	    mq.WithExchange("notifications"),
//	    mq.WithRoutingKey("candidate.rejected"),
//	)
type MQPublisher struct {
	log    *zap.Logger
	rabbit *RabbitMQ

	pub *rabbitmq.Publisher
}

var _ svc.Service = (*MQPublisher)(nil)

func NewMQPublisher(log *zap.Logger, rabbit *RabbitMQ) *MQPublisher {
	return &MQPublisher{
		log:    log,
		rabbit: rabbit,
	}
}

func (m *MQPublisher) Name() string        { return "rabbitmq_publisher" }
func (m *MQPublisher) DependsOn() []string { return []string{"logger", "rabbitmq"} }

// Init создаёт внутренний Publisher go-rabbitmq.
// Publisher управляет одним AMQP-каналом с поддержкой автоматического
// переподключения при обрыве соединения.
func (m *MQPublisher) Init(_ context.Context) error {
	pub, err := rabbitmq.NewPublisher(
		m.rabbit.Conn,
		rabbitmq.WithPublisherOptionsLogger(newZapLogger(m.log)),
	)
	if err != nil {
		return fmt.Errorf("failed to create rabbitmq publisher: %w", err)
	}

	m.pub = pub

	m.log.Debug("rabbitmq publisher initialized")

	return nil
}

// HealthCheck проверяет, что Publisher был инициализирован.
func (m *MQPublisher) HealthCheck(_ context.Context) error {
	if m.pub == nil {
		return errors.New("publisher is not initialized")
	}

	return nil
}

// Run — no-op: Publisher не имеет собственного event-loop.
// Публикация происходит синхронно по вызову Publish().
func (m *MQPublisher) Run(_ context.Context) error {
	return nil
}

// Stop корректно закрывает Publisher, завершая все ожидающие подтверждения (confirms).
func (m *MQPublisher) Stop(_ context.Context) error {
	m.pub.Close()

	return nil
}

// Publish публикует сообщение body в RabbitMQ.
//
// Параметры публикации (exchange, routing key, content type и т.д.)
// задаются через функциональные опции PublisherOption:
//
//	err := publisher.Publish(ctx, jsonBytes,
//	    mq.WithExchange("hiring.events"),
//	    mq.WithRoutingKey("candidate.created"),
//	)
//
// Если exchange или routing key не заданы явно — используются пустые значения
// (default exchange в RabbitMQ — прямая доставка по имени очереди).
func (m *MQPublisher) Publish(ctx context.Context, body []byte, opts ...PublisherOption) error {
	// Применяем все переданные опции к дефолтным настройкам.
	o := defaultPublishOptions()
	for _, opt := range opts {
		opt(&o)
	}

	pubOpts := []func(*rabbitmq.PublishOptions){
		rabbitmq.WithPublishOptionsExchange(o.exchange),
		rabbitmq.WithPublishOptionsContentType(o.contentType),
		rabbitmq.WithPublishOptionsPersistentDelivery, // сообщения сохраняются на диск
	}

	// WithPublishOptionsMandatory — это не factory-функция, а прямой setter.
	// Добавляем её только если mandatory = true, чтобы избежать лишних аллокаций.
	if o.mandatory {
		pubOpts = append(pubOpts, rabbitmq.WithPublishOptionsMandatory)
	}

	return m.pub.PublishWithContext(
		ctx,
		body,
		[]string{o.routingKey},
		pubOpts...,
	)
}
