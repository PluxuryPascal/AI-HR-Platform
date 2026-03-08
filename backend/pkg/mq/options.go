package mq

// ═══════════════════════════════════════════════════════════════════════════
// PUBLISHER OPTIONS
// ═══════════════════════════════════════════════════════════════════════════

// publishOptions хранит параметры одного вызова Publish.
// Значения по умолчанию заданы в defaultPublishOptions().
type publishOptions struct {
	// exchange — имя обменника (Exchange), куда публикуется сообщение.
	//
	// Exchange — это «маршрутизатор» RabbitMQ. Publisher никогда не пишет
	// напрямую в очередь. Он публикует в Exchange, а тот по типу
	// (topic/direct/fanout) и routing key решает, в какие очереди
	// переложить сообщение.
	//
	// Если пуст — используется default exchange (доставка по имени очереди).
	exchange string

	// routingKey — ключ маршрутизации.
	// Для topic-exchange: шаблон "candidate.created", "candidate.#".
	// Для direct-exchange: точное совпадение строки.
	// Для fanout-exchange: игнорируется (доставка всем подписчикам).
	routingKey string

	// mandatory — если true и ни одна очередь не привязана к routing key,
	// RabbitMQ вернёт сообщение обратно в Return-обработчик.
	// Обычно false (silently drop необработанные сообщения).
	mandatory bool

	// contentType — MIME-тип тела сообщения.
	// По умолчанию "application/json".
	contentType string
}

func defaultPublishOptions() publishOptions {
	return publishOptions{contentType: "application/json"}
}

// PublisherOption — функциональная опция для одного вызова Publish.
//
// Паттерн Functional Options: вместо передачи struct с кучей полей
// (большинство из которых всегда nil/zero) — передаём только нужные
// функции-опции. Читаемо, расширяемо, без brittle API.
//
// Пример — один publisher для нескольких exchange:
//
//	publisher.Publish(ctx, body, mq.WithExchange("hiring.events"), mq.WithRoutingKey("candidate.created"))
//	publisher.Publish(ctx, body, mq.WithExchange("notifications"), mq.WithRoutingKey("candidate.rejected"))
type PublisherOption func(*publishOptions)

// WithExchange задаёт имя Exchange, куда будет опубликовано сообщение.
func WithExchange(name string) PublisherOption {
	return func(o *publishOptions) { o.exchange = name }
}

// WithRoutingKey задаёт routing key для маршрутизации внутри Exchange.
func WithRoutingKey(key string) PublisherOption {
	return func(o *publishOptions) { o.routingKey = key }
}

// WithContentType задаёт MIME-тип тела сообщения (по умолчанию "application/json").
func WithContentType(ct string) PublisherOption {
	return func(o *publishOptions) { o.contentType = ct }
}

// WithMandatory включает флаг mandatory: RabbitMQ вернёт сообщение,
// если не найдётся ни одной очереди по routing key.
func WithMandatory() PublisherOption {
	return func(o *publishOptions) { o.mandatory = true }
}

// ═══════════════════════════════════════════════════════════════════════════
// CONSUMER OPTIONS
// ═══════════════════════════════════════════════════════════════════════════

// consumerConfig хранит всю конфигурацию Consumer.
//
// Ключевая идея: Consumer «владеет» своей очередью и объявляет её
// параметры самостоятельно при Init(). Это исключает необходимость
// в централизованном topology.go и дает каждому Consumer гибкость.
type consumerConfig struct {
	// queueName — имя очереди, которую слушает этот Consumer.
	// Пример: "hiring.candidate.created".
	queueName string

	// exchange — имя Exchange, к которому привязывается очередь.
	// Пример: "hiring.events".
	exchange string

	// routingKey — ключ привязки (binding) очереди к Exchange.
	// Пример: "candidate.created".
	routingKey string

	// consumerName — уникальный идентификатор consumer-а на стороне брокера.
	// Отображается в RabbitMQ Management UI. Если пуст — генерируется автоматически.
	consumerName string

	// prefetchCount — количество «в-полёте» (unacked) сообщений.
	//
	// При prefetch=1: строгая последовательность, низкий throughput.
	// При prefetch=10: до 10 сообщений параллельно, высокий throughput.
	//
	// Для тяжёлых операций (OCR, LLM) рекомендуется prefetch=1.
	prefetchCount int

	// quorum — объявить очередь как Quorum Queue (x-queue-type: quorum).
	//
	// Quorum Queue vs Classic Queue:
	//   Classic: данные только на одной ноде → потеря при падении ноды.
	//   Quorum: Raft-репликация на N нод (≥3) → выживает при падении (N-1)/2 нод.
	//
	// В production всегда использовать Quorum Queues для critical data.
	// Нельзя конвертировать существующую Classic в Quorum — нужно пересоздать.
	quorum bool

	// dlxExchange — Dead Letter Exchange (DLX).
	//
	// Если задан, RabbitMQ автоматически перемещает сообщения сюда при:
	//   1. NackDiscard от consumer-а (постоянная ошибка, requeue=false).
	//   2. Истечении TTL (x-message-ttl).
	//   3. Превышении лимита доставок (x-delivery-limit, только Quorum).
	//
	// DLX — это просто ещё один Exchange (обычно direct-type).
	// Отдельный Consumer (DLQ Worker) подписывается на Dead Letter Queue
	// и обрабатывает «мёртвые» сообщения (алерты, gRPC FAILED callback).
	dlxExchange string

	// dlxRoutingKey — routing key для сообщений в DLX.
	// Он же — имя Dead Letter Queue (DLQ).
	// Пример: "hiring.candidate.dead".
	dlxRoutingKey string

	// messageTTL — максимальное время жизни сообщения в очереди (миллисекунды).
	//
	// Если AI Engine не забрал сообщение за это время (например, сервис лежит),
	// RabbitMQ перекладывает его в DLX. Hiring Service затем пометит кандидата
	// как FAILED через DLQ Worker.
	//
	// Пример: 300_000 = 5 минут.
	messageTTL int

	// maxDeliveries — x-delivery-limit (только для Quorum Queues).
	//
	// Если handler возвращает NackRequeue N раз подряд (например, сеть падает
	// при каждой попытке), после N-й попытки RabbitMQ перестаёт requeue
	// и отправляет сообщение в DLX автоматически.
	//
	// Предотвращает бесконечный цикл: получить → ошибка → вернуть → повторить...
	maxDeliveries int

	// concurrency — количество горутин, обрабатывающих сообщения параллельно
	// внутри одного инстанса Consumer-а. По умолчанию 1 (строго последовательно).
	// В сочетании с prefetchCount позволяет балансировать throughput.
	concurrency int
}

func defaultConsumerConfig() consumerConfig {
	return consumerConfig{
		prefetchCount: 10,
		quorum:        true, // quorum по умолчанию — безопасный default для production
		concurrency:   1,    // по умолчанию 1 обработчик на Consumer (последовательно)
	}
}

// ConsumerOption — функциональная опция для настройки Consumer.
// Применяется в NewMQConsumer.
type ConsumerOption func(*consumerConfig)

// WithQueueName задаёт имя очереди, которую будет слушать Consumer.
func WithQueueName(name string) ConsumerOption {
	return func(c *consumerConfig) { c.queueName = name }
}

// WithConsumerExchange задаёт Exchange, к которому привязывается очередь.
func WithConsumerExchange(name string) ConsumerOption {
	return func(c *consumerConfig) { c.exchange = name }
}

// WithConsumerRoutingKey задаёт routing key для привязки очереди к Exchange.
func WithConsumerRoutingKey(key string) ConsumerOption {
	return func(c *consumerConfig) { c.routingKey = key }
}

// WithConsumerName задаёт имя consumer-а в Management UI RabbitMQ.
func WithConsumerName(name string) ConsumerOption {
	return func(c *consumerConfig) { c.consumerName = name }
}

// WithPrefetchCount задаёт количество unacked сообщений (QoS prefetch).
// При prefetch=1 — строгая последовательная обработка.
func WithPrefetchCount(n int) ConsumerOption {
	return func(c *consumerConfig) { c.prefetchCount = n }
}

// WithConcurrency задаёт количество горутин (воркеров), которые будут
// параллельно вызывать handler для полученных сообщений ресурса Consumer.
func WithConcurrency(n int) ConsumerOption {
	return func(c *consumerConfig) { c.concurrency = n }
}

// WithQuorumQueue объявляет очередь как Quorum Queue (x-queue-type: quorum).
//
// Quorum Queues реплицируются между нодами кластера через Raft-протокол.
// Потеря одной ноды не ведёт к потере сообщений.
// Рекомендуется для production (включён по умолчанию в defaultConsumerConfig).
func WithQuorumQueue() ConsumerOption {
	return func(c *consumerConfig) { c.quorum = true }
}

// WithDLX настраивает Dead Letter Exchange для очереди.
//
// dlxExchange — имя обменника для «мёртвых» сообщений (например, "hiring.dlx").
// dlxRoutingKey — routing key (= имя DLQ) (например, "hiring.candidate.dead").
//
// Важно: DLQ (Dead Letter Queue) объявляется не здесь, а отдельным Consumer-ом
// (DLQ Worker), который подписывается на dlqName ("hiring.candidate.dead").
// Это гарантирует существование DLQ без сырого amqp091-go кода.
func WithDLX(dlxExchange, dlxRoutingKey string) ConsumerOption {
	return func(c *consumerConfig) {
		c.dlxExchange = dlxExchange
		c.dlxRoutingKey = dlxRoutingKey
	}
}

// WithMessageTTL задаёт время жизни сообщения в очереди (миллисекунды).
// По истечении TTL сообщение уходит в DLX.
// Пример: WithMessageTTL(300_000) = 5 минут.
func WithMessageTTL(ms int) ConsumerOption {
	return func(c *consumerConfig) { c.messageTTL = ms }
}

// WithMaxDeliveries задаёт x-delivery-limit (только для Quorum Queues).
// После N неудачных попыток доставки сообщение уходит в DLX.
// Работает совместно с NackRequeue для ограничения числа retry.
func WithMaxDeliveries(n int) ConsumerOption {
	return func(c *consumerConfig) { c.maxDeliveries = n }
}
