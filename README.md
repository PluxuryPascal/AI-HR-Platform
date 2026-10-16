# AI-HR-Platform

ИИ-платформа для найма: загружаете резюме кандидатов, а система сама разбирает их, оценивает под вакансию, сравнивает между собой и готовит письма кандидатам.

**Бэкенд:** микросервисы на Go, которые общаются по gRPC, длительные операции идут через Temporal, события через RabbitMQ. **Фронтенд:** Next.js.

## Возможности

- **Парсинг резюме** (PDF) и описаний вакансий с помощью LLM
- **Скоринг кандидатов** под вакансию с объяснением оценки
- **Семантический поиск** по резюме: эмбеддинги в PostgreSQL + pgvector
- **Сравнение кандидатов** и генерация персональных писем
- **Команды и роли**: мультитенантный RBAC на Casbin (роли действуют в рамках команды)
- **Уведомления** по email

## Архитектура

Сервисы общаются по gRPC (контракты в `proto/`, генерация через Buf). Асинхронные задачи идут через RabbitMQ, а длительные многошаговые процессы через Temporal.

| Сервис | Что делает |
|---|---|
| `auth` | пользователи, команды, приглашения, JWT, RBAC (Casbin) |
| `hiring` | вакансии, кандидаты, отделы, загрузка резюме |
| `ai-engine` | gRPC API для скоринга, сравнения, писем, разбора вакансий |
| `ai-worker` | Temporal-воркер: пайплайн обработки резюме, интервью, чат, DLQ-колбэки |
| `notification` | доставка email-уведомлений из очереди |

**Temporal-воркфлоу:** `resume_pipeline`, `job_parse`, `candidate_compare`, `email_generate`, `interview_workflow`, `chat_workflow`, `dlq_callback`. Ретраи и восстановление после сбоев берёт на себя Temporal, а не самописная логика.

## Стек

`Go` · `gRPC` + `Buf` · `Echo` · `Temporal` · `RabbitMQ` · `PostgreSQL` · `pgvector` · `Redis` · `Casbin` · `OpenRouter (LLM)` · `Cloudinary` · `Docker Compose` · `Next.js`

## Структура

```
backend/
├── cmd/            # точки входа сервисов: auth, hiring, ai_engine, ai_worker, notification
├── internal/
│   ├── domain/     # доменные модели
│   ├── usecase/    # бизнес-логика
│   ├── repo/       # доступ к PostgreSQL
│   ├── grpc-handler/, http-handler/
│   ├── temporal/   # workflow + activity
│   ├── llm/        # парсинг, скоринг, эмбеддинги, генерация писем
│   └── worker/     # консьюмеры RabbitMQ
├── proto/          # gRPC-контракты (Buf)
├── migrations/     # SQL-миграции
├── casbin/         # модель и политики RBAC
└── docker-compose.yml
web/                # Next.js-клиент
```

## Запуск

```bash
cd backend
docker compose up -d --build
```

Поднимутся PostgreSQL, Redis, RabbitMQ, Temporal (с Temporal UI), все сервисы и веб-клиент. Для ИИ-функций нужен ключ LLM-провайдера (OpenRouter) в переменных окружения.
