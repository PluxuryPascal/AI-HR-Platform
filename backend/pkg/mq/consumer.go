package mq

import (
	"backend/pkg/svc"
	"context"
	"errors"
	"fmt"

	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

// HandlerFunc — функция-обработчик одного входящего сообщения.
//
// # Параметры
//
//   - ctx — контекст вызова. Создаётся в MQConsumer.Run() и отменяется при
//     остановке сервиса. ОБЯЗАТЕЛЕН для:
//   - Graceful Shutdown: if ctx.Err() != nil — не запускать тяжёлую работу.
//   - Request Timeout: context.WithTimeout(ctx, 30*time.Second) для LLM-вызова.
//   - OpenTelemetry: извлечение trace ID из d.Headers и propagation в ctx.
//   - d — доставленное сообщение. Содержит Body (JSON/proto), Headers,
//     RoutingKey, DeliveryTag и другие метаданные.
//
// # Возвращаемые значения (rabbitmq.Action)
//
// Действие, которое Consumer автоматически применит к сообщению:
//
//   - rabbitmq.Ack — сообщение успешно обработано. Брокер удаляет его из очереди. ✅
//
//   - rabbitmq.NackRequeue — временная ошибка (сеть упала, S3 вернул 502).
//     Брокер возвращает сообщение в очередь для повторной попытки. 🔄
//     Quorum Queue отслеживает x-delivery-count. Когда count > x-delivery-limit,
//     RabbitMQ автоматически отправляет сообщение в DLX (без вашего участия). ☠️
//
//   - rabbitmq.NackDiscard — постоянная ошибка (PDF зашифрован паролем, S3 404).
//     Брокер НЕ возвращает сообщение в очередь. Если настроен DLX —
//     сообщение попадёт туда. DLQ Worker вызовет gRPC PARSING_STATUS_FAILED. ☠️
type HandlerFunc func(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action

// MQConsumer реализует svc.Service и управляет одним Consumer-ом RabbitMQ.
//
// # Ключевые принципы
//
//  1. Handler инжектируется — Consumer не знает бизнес-логику.
//     AI Engine передаёт свой хэндлер, DLQ Worker — свой.
//
//  2. Consumer сам объявляет свою очередь при Init() через QueueArgs.
//     Нет централизованного topology.go — каждый Consumer «владеет» очередью.
//
//  3. Run() блокируется до закрытия consumer-а.
//     svc.Run() запускает Run() в goroutine и ждёт её завершения.
//
//  4. Manual Ack (autoAck=false) — сообщение считается обработанным
//     только после явного Ack/Nack от хэндлера. Сбой процесса → requeue.
type MQConsumer struct {
	log      *zap.Logger
	rabbitMQ *RabbitMQ
	handler  HandlerFunc
	cfg      consumerConfig
	consumer *rabbitmq.Consumer
}

var _ svc.Service = (*MQConsumer)(nil)

// NewMQConsumer создаёт Consumer с инжектируемым хэндлером.
//
// Пример — основной Consumer AI Engine:
//
//	mq.NewMQConsumer(log, rmq.Conn, aiHandler,
//	    mq.WithQueueName("hiring.candidate.created"),
//	    mq.WithConsumerExchange("hiring.events"),
//	    mq.WithConsumerRoutingKey("candidate.created"),
//	    mq.WithQuorumQueue(),
//	    mq.WithDLX("hiring.dlx", "hiring.candidate.dead"),
//	    mq.WithMessageTTL(300_000),
//	    mq.WithMaxDeliveries(3),
//	    mq.WithPrefetchCount(1), // AI обрабатывает по одному (тяжёлые задачи)
//	)
//
// Пример — DLQ Worker (объявляет DLQ-очередь, тем самым создавая её):
//
//	mq.NewMQConsumer(log, rmq.Conn, dlqHandler,
//	    mq.WithQueueName("hiring.candidate.dead"),
//	    mq.WithConsumerExchange("hiring.dlx"),
//	    mq.WithConsumerRoutingKey("hiring.candidate.dead"),
//	)
func NewMQConsumer(
	log *zap.Logger,
	rabbitMQ *RabbitMQ,
	handler HandlerFunc,
	opts ...ConsumerOption,
) *MQConsumer {
	cfg := defaultConsumerConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	return &MQConsumer{
		log:      log,
		rabbitMQ: rabbitMQ,
		handler:  handler,
		cfg:      cfg,
	}
}

func (m *MQConsumer) Name() string {
	if m.cfg.serviceName != "" {
		return m.cfg.serviceName
	}
	return "rabbitmq_consumer"
}
func (m *MQConsumer) DependsOn() []string { return []string{"logger", "rabbitmq"} }

// Init создаёт внутренний Consumer и объявляет очередь.
//
// При вызове rabbitmq.NewConsumer библиотека:
//  1. Декларирует Exchange (если указан WithConsumerOptionsExchangeDeclare).
//  2. Декларирует Queue с переданными QueueArgs (quorum, DLX, TTL...).
//  3. Создаёт binding: Queue ↔ Exchange + RoutingKey.
//  4. Настраивает QoS (prefetch).
//
// Это ключевое отличие от подхода с topology.go: каждый Consumer сам
// объявляет нужную ему инфраструктуру. Нет «централизованного» места,
// которое нужно синхронизировать с Consumer-ами.
func (m *MQConsumer) Init(_ context.Context) error {
	if m.rabbitMQ.Conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	consumer, err := rabbitmq.NewConsumer(
		m.rabbitMQ.Conn,
		m.cfg.queueName,
		m.buildOptions()...,
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer for queue %q: %w", m.cfg.queueName, err)
	}

	m.consumer = consumer

	m.log.Debug("rabbitmq consumer initialized",
		zap.String("queue", m.cfg.queueName),
		zap.String("exchange", m.cfg.exchange),
		zap.String("routing_key", m.cfg.routingKey),
		zap.Bool("quorum", m.cfg.quorum),
		zap.String("dlx", m.cfg.dlxExchange),
	)

	return nil
}

// HealthCheck проверяет инициализацию Consumer.
func (m *MQConsumer) HealthCheck(_ context.Context) error {
	if m.consumer == nil {
		return errors.New("consumer is not initialized")
	}

	return nil
}

// Run запускает блокирующий event-loop получения и обработки сообщений.
//
// consumer.Run() из go-rabbitmq блокируется до закрытия consumer-а.
// Передаём ctx в handler — это обязательно для graceful shutdown,
// таймаутов и OpenTelemetry propagation.
//
// go-rabbitmq автоматически применяет Ack/Nack по возвращаемому Action.
func (m *MQConsumer) Run(ctx context.Context) error {
	if err := m.consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		return m.handler(ctx, d)
	}); err != nil {
		return fmt.Errorf("consumer %q stopped unexpectedly: %w", m.cfg.queueName, err)
	}

	return nil
}

// Stop корректно закрывает Consumer, дожидаясь завершения текущего handler-а
// (go-rabbitmq по умолчанию ждёт handler с CloseGracefully=true).
func (m *MQConsumer) Stop(_ context.Context) error {
	m.consumer.Close()

	return nil
}

// buildOptions конвертирует consumerConfig в slice опций go-rabbitmq.
func (m *MQConsumer) buildOptions() []func(*rabbitmq.ConsumerOptions) {
	// queueArgs — «аргументы» AMQP-протокола при объявлении очереди.
	// Это key-value таблица, которая передаётся брокеру.
	// Здесь описываются все «расширенные» характеристики очереди.
	queueArgs := rabbitmq.Table{}

	if m.cfg.dlxExchange != "" {
		// x-dead-letter-exchange: при Nack/TTL — куда перекладывать сообщения
		queueArgs["x-dead-letter-exchange"] = m.cfg.dlxExchange
	}

	if m.cfg.dlxRoutingKey != "" {
		// x-dead-letter-routing-key: с каким ключом перекладывать в DLX
		queueArgs["x-dead-letter-routing-key"] = m.cfg.dlxRoutingKey
	}

	if m.cfg.messageTTL > 0 {
		// x-message-ttl: TTL сообщения в миллисекундах
		queueArgs["x-message-ttl"] = m.cfg.messageTTL
	}

	if m.cfg.maxDeliveries > 0 {
		// x-delivery-limit: макс. попыток для Quorum Queue
		queueArgs["x-delivery-limit"] = m.cfg.maxDeliveries
	}

	opts := []func(*rabbitmq.ConsumerOptions){
		// Manual Ack: брокер ждёт явного Ack от handler-а.
		rabbitmq.WithConsumerOptionsConsumerAutoAck(false),
		// Привязка к Exchange.
		rabbitmq.WithConsumerOptionsExchangeName(m.cfg.exchange),
	}

	if m.cfg.declareExchange {
		opts = append(opts, rabbitmq.WithConsumerOptionsExchangeDeclare)
		if m.cfg.exchangeType != "" {
			opts = append(opts, rabbitmq.WithConsumerOptionsExchangeKind(m.cfg.exchangeType))
		} else {
			opts = append(opts, rabbitmq.WithConsumerOptionsExchangeKind("topic"))
		}
		opts = append(opts, rabbitmq.WithConsumerOptionsExchangeDurable)
	}

	opts = append(opts,
		// Routing key для binding Queue ↔ Exchange.
		rabbitmq.WithConsumerOptionsRoutingKey(m.cfg.routingKey),
		// QueueArgs с DLX, TTL, delivery-limit.
		rabbitmq.WithConsumerOptionsQueueArgs(queueArgs),
		// QoS: prefetch count.
		rabbitmq.WithConsumerOptionsQOSPrefetch(m.cfg.prefetchCount),
		// Durable: очередь переживёт перезапуск брокера.
		rabbitmq.WithConsumerOptionsQueueDurable,
		// Логгер.
		rabbitmq.WithConsumerOptionsLogger(newZapLogger(m.log)),
	)

	if m.cfg.quorum {
		// Нативный метод библиотеки для x-queue-type: quorum.
		// Нельзя задать вручную через queueArgs, если очередь уже существует
		// как classic — в этом случае нужно переименовать/пересоздать.
		opts = append(opts, rabbitmq.WithConsumerOptionsQueueQuorum)
	}

	if m.cfg.consumerName != "" {
		opts = append(opts, rabbitmq.WithConsumerOptionsConsumerName(m.cfg.consumerName))
	}

	// Устанавливаем количество параллельных горутин-воркеров для этого Consumer-а.
	if m.cfg.concurrency > 1 {
		opts = append(opts, rabbitmq.WithConsumerOptionsConcurrency(m.cfg.concurrency))
	}

	return opts
}
