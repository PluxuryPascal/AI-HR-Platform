// Package mq — инфраструктурный пакет для работы с RabbitMQ.
//
// Этот файл содержит конкретный пример сборки Consumer-ов для сервиса AI Engine
// (обработка резюме) и DLQ Worker-а (обработка «мёртвых» сообщений).
//
// Он НЕ является частью пакета mq — он находится здесь как живая документация.
// В реальном приложении эти функции/структуры размещаются в соответствующих
// пакетах (например, cmd/ai-engine/main.go или internal/worker/candidate.go).
package mq

// ═══════════════════════════════════════════════════════════════════════════
// ПРИМЕР 1: Сборка Consumer-ов в cmd/ai-engine/main.go
// ═══════════════════════════════════════════════════════════════════════════

// ExampleSetup показывает, как собрать основной Consumer и DLQ Worker.
//
// В реальном приложении этот код находится в main() или wire.go.
//
//	func ExampleSetup() {
//	    log, _ := zap.NewProduction()
//	    rmq := mq.NewRabbitMQ(log, &cfg.RabbitMQ)
//
//	    // ── Основной Consumer (AI Engine) ────────────────────────────────
//	    // Слушает "hiring.candidate.created", обрабатывает резюме.
//	    // Если AI не смог обработать за 5 минут или > 3 попыток → DLX.
//	    mainConsumer := mq.NewMQConsumer(
//	        log,
//	        rmq.Conn,
//	        CandidateCreatedHandler,  // бизнес-логика AI Engine (см. ниже)
//	        mq.WithQueueName("hiring.candidate.created"),
//	        mq.WithConsumerExchange("hiring.events"),
//	        mq.WithConsumerRoutingKey("candidate.created"),
//	        mq.WithQuorumQueue(),                               // Quorum: репликация
//	        mq.WithDLX("hiring.dlx", "hiring.candidate.dead"), // DLX config
//	        mq.WithMessageTTL(300_000),                        // 5 минут TTL
//	        mq.WithMaxDeliveries(3),                           // макс. 3 попытки NackRequeue
//	        mq.WithPrefetchCount(20),                          // берем пачку в 20 сообщений
//	        mq.WithConcurrency(10),                            // 10 горутин параллельно (LLM I/O bound)
//	        mq.WithConsumerName("ai-engine-candidate-processor"),
//	    )
//
//	    // ── DLQ Worker ───────────────────────────────────────────────────
//	    // Слушает "hiring.candidate.dead" — очередь, куда RabbitMQ складывает
//	    // «мёртвые» сообщения из основной очереди.
//	    //
//	    // КЛЮЧЕВОЙ МОМЕНТ: создавая Consumer для "hiring.candidate.dead",
//	    // мы тем самым ОБЪЯВЛЯЕМ эту очередь в RabbitMQ.
//	    // Никакого topology.go не нужно — Consumer объявляет очередь сам.
//	    //
//	    // DLQ Worker привязывается к Exchange "hiring.dlx"
//	    // с routing key "hiring.candidate.dead".
//	    // Когда mainConsumer делает NackDiscard или не обрабатывает за 5 мин,
//	    // RabbitMQ перекладывает сообщение в hiring.dlx → hiring.candidate.dead.
//	    dlqWorker := mq.NewMQConsumer(
//	        log,
//	        rmq.Conn,
//	        DLQCandidateHandler,    // отдельный хэндлер для «мёртвых» сообщений
//	        mq.WithQueueName("hiring.candidate.dead"),
//	        mq.WithConsumerExchange("hiring.dlx"),
//	        mq.WithConsumerRoutingKey("hiring.candidate.dead"),
//	        // DLQ — classic queue (не quorum): для ручного разбора,
//	        // не нужна репликация. Это обычно небольшой объём сообщений.
//	        mq.WithConsumerName("ai-engine-dlq-worker"),
//	    )
//
//	    // Регистрируем в svc.Run — оба Consumer-а запускаются параллельно.
//	    svc.Run(ctx, log, []svc.Service{
//	        logSvc,
//	        rmq,          // DependsOn: ["logger"]
//	        mainConsumer, // DependsOn: ["logger", "rabbitmq"]
//	        dlqWorker,    // DependsOn: ["logger", "rabbitmq"]
//	    })
//	}

// ═══════════════════════════════════════════════════════════════════════════
// ПРИМЕР 2: Основной хэндлер с Human-in-the-Loop логикой
// ═══════════════════════════════════════════════════════════════════════════

// CandidateCreatedHandler — хэндлер основного Consumer-а.
// Вызывается при получении события "candidate.created".
//
// # Human-in-the-Loop (HIL) Algorithm
//
// Наш pipeline обработки резюме имеет три исхода:
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │ 1. ВРЕМЕННАЯ ИНФРА-ОШИБКА (сеть, S3 502)                           │
// │    Handler → NackRequeue                                            │
// │    Quorum Queue: x-delivery-count++                                 │
// │    Если count > x-delivery-limit → автоматически в DLX/DLQ         │
// ├─────────────────────────────────────────────────────────────────────┤
// │ 2. ПОСТОЯННАЯ ИНФРА-ОШИБКА (PDF зашифрован, S3 404)                │
// │    Handler → NackDiscard                                            │
// │    RabbitMQ → DLX → DLQ ("hiring.candidate.dead")                  │
// │    DLQ Worker → gRPC UpdateCandidateProfile(PARSING_STATUS_FAILED) │
// ├─────────────────────────────────────────────────────────────────────┤
// │ 3. БИЗНЕС-СЦЕНАРИЙ — Human-in-the-Loop                             │
// │    AI нашёл skills/experience, но не нашёл email или first_name.   │
// │    Handler → gRPC UpdateCandidateProfile(NEEDS_REVIEW,             │
// │                  structured_data={skills, experience})             │
// │    Handler → Ack  ← ключевое отличие от п.2!                       │
// │    Сообщение УДАЛЕНО из очереди. Фронтенд предложит рекрутеру      │
// │    заполнить email вручную (Smart Autofill modal).                  │
// └─────────────────────────────────────────────────────────────────────┘
//
//	func CandidateCreatedHandler(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action {
//	    // Сначала проверяем контекст: если ctx отменён (graceful shutdown),
//	    // не начинаем тяжёлую работу — вернём в очередь для следующей попытки.
//	    if ctx.Err() != nil {
//	        return rabbitmq.NackRequeue
//	    }
//
//	    // Десериализуем событие из JSON.
//	    var event hiringpb.CandidateCreatedEvent
//	    if err := protojson.Unmarshal(d.Body, &event); err != nil {
//	        // Битый JSON — постоянная ошибка. Повтор не поможет.
//	        // → NackDiscard → DLX → DLQ Worker сообщит Hiring о FAILED.
//	        log.Error("failed to unmarshal event", zap.Error(err))
//	        return rabbitmq.NackDiscard
//	    }
//
//	    // Скачиваем файл из S3.
//	    fileBytes, err := s3Client.Download(ctx, event.ResumeFileKey)
//	    if err != nil {
//	        if isTemporary(err) {
//	            // Сеть упала, S3 временно недоступен — попробуем снова.
//	            log.Warn("temporary S3 error, requeuing", zap.Error(err))
//	            return rabbitmq.NackRequeue // 🔄 повтор
//	        }
//	        // S3 404 или другой постоянный сбой — файл никогда не появится.
//	        log.Error("permanent S3 error", zap.Error(err))
//	        return rabbitmq.NackDiscard // ☠️ → DLX
//	    }
//
//	    // Детектируем тип файла.
//	    if isPasswordProtected(fileBytes) {
//	        // PDF зашифрован паролем — AI не может распарсить.
//	        // Это постоянная ошибка — повтор не поможет.
//	        log.Error("password-protected PDF, discarding", zap.String("key", event.ResumeFileKey))
//	        return rabbitmq.NackDiscard // ☠️ → DLX → DLQ Worker → gRPC FAILED
//	    }
//
//	    // Запускаем OCR + LLM. Даём жёсткий таймаут (LLM может зависнуть).
//	    parseCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
//	    defer cancel()
//
//	    result, err := aiEngine.ParseResume(parseCtx, fileBytes)
//	    if err != nil {
//	        if errors.Is(err, context.DeadlineExceeded) {
//	            // LLM завис — временная ошибка, попробуем ещё раз.
//	            log.Warn("LLM timeout, requeuing")
//	            return rabbitmq.NackRequeue // 🔄
//	        }
//	        // Неизвестная ошибка — в DLQ для ручного разбора.
//	        log.Error("unexpected parse error", zap.Error(err))
//	        return rabbitmq.NackDiscard // ☠️
//	    }
//
//	    // ─────────────────────────────────────────────────────────────────
//	    // HUMAN-IN-THE-LOOP: AI успешно извлёк данные, но часть отсутствует.
//	    // ─────────────────────────────────────────────────────────────────
//	    //
//	    // Сценарий: AI нашёл опыт работы, навыки, образование —
//	    // но не смог извлечь email (нечёткое форматирование PDF).
//	    //
//	    // ЧТО МЫ ДЕЛАЕМ:
//	    //   1. Вызываем gRPC UpdateCandidateProfile с NEEDS_REVIEW.
//	    //      Передаём ВСЁ, что AI нашёл (skills, experience_years...).
//	    //   2. Возвращаем Ack — сообщение удалено из очереди.
//	    //
//	    // ЧТО ПРОИСХОДИТ ДАЛЬШЕ:
//	    //   Hiring Service сохраняет structured_data, помечает кандидата.
//	    //   Фронтенд подсвечивает профиль → рекрутер вводит email вручную.
//	    //   Это Human-in-the-Loop: ценные данные сохранены, человек
//	    //   дополняет только отсутствующие поля.
//	    //
//	    // ПОЧЕМУ НЕ NackDiscard?
//	    //   NackDiscard означает «инфраструктурная ошибка, данные утеряны».
//	    //   Здесь AI сработал нормально — просто данных недостаточно.
//	    //   Конечный статус NEEDS_REVIEW — это бизнес-решение, не ошибка брокера.
//
//	    status := hiringpb.ParsingStatus_PARSING_STATUS_SUCCESS
//	    if result.RequiresHumanReview() { // нет email, first_name и т.п.
//	        status = hiringpb.ParsingStatus_PARSING_STATUS_NEEDS_REVIEW
//	    }
//
//	    // Вызываем gRPC callback в Hiring Service.
//	    _, err = hiringClient.UpdateCandidateProfile(ctx, &hiringpb.UpdateCandidateProfileRequest{
//	        CandidateId:    event.CandidateId,
//	        ParsingStatus:  status,
//	        ParsedText:     result.RawText,
//	        StructuredData: result.ToProtoStruct(), // skills, experience, education...
//	    })
//	    if err != nil {
//	        // gRPC временно недоступен — попробуем снова (весь pipeline).
//	        log.Warn("grpc callback failed, requeuing", zap.Error(err))
//	        return rabbitmq.NackRequeue // 🔄
//	    }
//
//	    // Всё сделано: gRPC вызван, данные сохранены, очередь очищена.
//	    return rabbitmq.Ack // ✅
//	}

// ═══════════════════════════════════════════════════════════════════════════
// ПРИМЕР 3: DLQ Worker хэндлер
// ═══════════════════════════════════════════════════════════════════════════

// DLQCandidateHandler — хэндлер для Dead Letter Queue.
// Вызывается, когда основной Consumer не смог обработать сообщение:
//   - NackDiscard (постоянная ошибка)
//   - Истёк TTL (AI Engine лежал более 5 минут)
//   - Превышен x-delivery-limit (3 NackRequeue подряд)
//
// Задача: сообщить Hiring Service о финальном FAILED статусе,
// чтобы рекрутер увидел, что кандидат не был обработан.
//
//	func DLQCandidateHandler(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action {
//	    // RabbitMQ добавляет x-death headers с метаданными о смерти сообщения:
//	    // - причина (rejected/expired/maxlen)
//	    // - исходная очередь
//	    // - timestamp
//	    // Их можно логировать для observability.
//	    if death, ok := d.Headers["x-death"]; ok {
//	        log.Warn("processing dead letter message", zap.Any("x-death", death))
//	    }
//
//	    var event hiringpb.CandidateCreatedEvent
//	    if err := protojson.Unmarshal(d.Body, &event); err != nil {
//	        // Если не можем даже десериализовать — выбрасываем из DLQ.
//	        // Дальнейший retry бессмысленен. Логируем для ручного разбора.
//	        log.Error("DLQ: cannot unmarshal dead letter, discarding permanently",
//	            zap.Error(err), zap.ByteString("body", d.Body))
//	        return rabbitmq.Ack // Ack, чтобы убрать из DLQ (уже залогировали)
//	    }
//
//	    // Вызываем gRPC с FAILED — Hiring Service пометит кандидата.
//	    // Фронтенд покажет рекрутеру статус "Не удалось обработать резюме".
//	    _, err := hiringClient.UpdateCandidateProfile(ctx, &hiringpb.UpdateCandidateProfileRequest{
//	        CandidateId:   event.CandidateId,
//	        ParsingStatus: hiringpb.ParsingStatus_PARSING_STATUS_FAILED,
//	    })
//	    if err != nil {
//	        // gRPC временно недоступен — оставляем в DLQ, попробуем позже.
//	        log.Error("DLQ: grpc callback failed", zap.Error(err))
//	        return rabbitmq.NackRequeue
//	    }
//
//	    log.Warn("DLQ: candidate marked as FAILED", zap.String("candidate_id", event.CandidateId))
//	    return rabbitmq.Ack // убираем из DLQ — Hiring оповещён
//	}

// ═══════════════════════════════════════════════════════════════════════════
// ПРИМЕР 4: Publisher из UseCase Hiring Service
// ═══════════════════════════════════════════════════════════════════════════

// ExamplePublish демонстрирует публикацию события из Hiring UseCase.
//
//	func (uc *createCandidateUC) Execute(ctx context.Context, req CreateCandidateRequest) error {
//	    // ... создаём кандидата в БД, загружаем файл в S3 ...
//
//	    payload, _ := json.Marshal(CandidateCreatedEvent{
//	        CandidateID:   candidate.ID.String(),
//	        ResumeFileKey: s3Key,
//	    })
//
//	    // Один publisher — один вызов с нужными параметрами.
//	    // Не нужно создавать отдельный Publisher для каждого topic.
//	    return uc.publisher.Publish(ctx, payload,
//	        mq.WithExchange("hiring.events"),
//	        mq.WithRoutingKey("candidate.created"),
//	    )
//	}

// ═══════════════════════════════════════════════════════════════════════════
// ДИАГРАММА ПОТОКА ДАННЫХ
// ═══════════════════════════════════════════════════════════════════════════

// Полный поток RabbitMQ для batch resume processing:
//
//	Hiring UseCase
//	  │ Publish("candidate.created")
//	  ▼
//	Exchange: "hiring.events" (topic)
//	  │ routing key: "candidate.created"
//	  ▼
//	Queue: "hiring.candidate.created" [Quorum, DLX=hiring.dlx, TTL=5min, MaxDeliv=3]
//	  │
//	  ├─── AI Engine Consumer (prefetch=1)
//	  │      │
//	  │      ├─ Временная ошибка → NackRequeue → [retry, до 3x] ──────┐
//	  │      │                                                          │ x-delivery-limit пройден
//	  │      ├─ Постоянная ошибка → NackDiscard ──────────────────────┤
//	  │      │                                                          │
//	  │      ├─ TTL истёк (AI Engine лежал) ──────────────────────────┤
//	  │      │                                                          ▼
//	  │      └─ Успех (SUCCESS или NEEDS_REVIEW) → gRPC → Ack     Exchange: "hiring.dlx" (direct)
//	  │                                                                 │ routing key: "hiring.candidate.dead"
//	  │                                                                 ▼
//	  │                                                        Queue: "hiring.candidate.dead" [DLQ]
//	  │                                                                 │
//	  │                                                        DLQ Worker Consumer
//	  │                                                                 │
//	  │                                                        gRPC UpdateCandidateProfile(FAILED)
//	  │                                                                 │
//	  │                                                                 └─ Ack
//	  │
//	  └─── Hiring gRPC Server
//	         │ UpdateCandidateProfile(SUCCESS/NEEDS_REVIEW/FAILED)
//	         ▼
//	       t_candidates + t_candidate_profiles (PostgreSQL)
//	         │
//	       WebSocket Hub → Frontend (real-time status update)
