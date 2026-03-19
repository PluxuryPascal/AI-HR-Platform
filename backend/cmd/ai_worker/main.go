package main

import (
	"backend/internal/audit"
	"backend/internal/db"
	"backend/internal/llm"
	"backend/internal/repo"
	"backend/internal/temporal"
	"backend/internal/temporal/activity"
	"backend/internal/worker"
	"backend/pkg/config"
	"backend/pkg/grpc"
	"backend/pkg/logger"
	"backend/pkg/mq"
	"backend/pkg/pdf"
	"backend/pkg/storage"
	"backend/pkg/svc"
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
)

type infrastructureComponents struct {
	cfg            *config.Config
	log            *logger.Log
	pool           *db.PostgresClient
	cloudinary     *storage.CloudinaryStorage
	llmClient      *llm.Client
	temporalClient *temporal.Client
	temporalWorker *temporal.Worker
	mqClient       *mq.RabbitMQ
	mqConsumer     *mq.MQConsumer
	dlqWorker      *mq.MQConsumer
	hiringGRPC     *grpc.Client
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("application error: %v", err)
	}

	log.Println("success shutdown")
}

func run(ctx context.Context) error {
	infra, err := initInfrastructure(ctx)
	if err != nil {
		return fmt.Errorf("init infrastructure error: %w", err)
	}

	if err := svc.Run(ctx, infra.log.Log, []svc.Service{
		infra.log,
		infra.pool,
		infra.cloudinary,
		infra.llmClient,
		infra.hiringGRPC,
		infra.temporalClient,
		infra.temporalWorker,
		infra.mqClient,
		infra.mqConsumer,
		infra.dlqWorker,
	}); err != nil {
		return fmt.Errorf("run service error: %w", err)
	}

	return nil
}

func initInfrastructure(ctx context.Context) (*infrastructureComponents, error) {
	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("load config error: %w", err)
	}

	zapLog, err := logger.New(
		logger.WithFile(conf.Logger.File.Path),
		logger.WithLevel(conf.Logger.Level),
		logger.WithStdOut(conf.Logger.StdOut),
	)
	if err != nil {
		return nil, fmt.Errorf("create logger error: %w", err)
	}

	pool, err := db.NewDb(zapLog.Log, conf)
	if err != nil {
		return nil, fmt.Errorf("create db error: %w", err)
	}

	cloudinary := storage.NewCloudinaryStorage(zapLog.Log, &conf.Cloudinary)
	llmClient := llm.NewClient(zapLog.Log, &conf.OpenRouter)
	temporalClient := temporal.NewClient(zapLog.Log, nil, &conf.Temporal)

	candidateRepo := repo.NewCandidateRepo(pool)
	commRepo := repo.NewCommunicationRepo(pool)
	jobRepo := repo.NewJobRepo(pool)

	hiringClientCfg := conf.GRPC.Clients["hiring"]
	hiringGRPC := grpc.NewClient("hiring", zapLog.Log, &hiringClientCfg)

	pdfExtractor := pdf.NewExtractor(pdf.Config{})

	auditor := audit.NewLogger(zapLog.Log, pool)
	if err := auditor.SeedActionTypes(ctx); err != nil {
		zapLog.Log.Warn("failed to seed audit action types", zap.Error(err))
	}

	aiSettingsRepo := repo.NewAiSettingsRepo(pool)
	llmProvider := llm.NewProvider(zapLog.Log, &conf.OpenRouter, aiSettingsRepo)

	activities := activity.NewActivities(
		zapLog.Log,
		pdfExtractor,
		cloudinary,
		llm.NewResumeParser(llmProvider),
		llm.NewScorer(llmProvider),
		llm.NewEmbedder(llmProvider),
		llm.NewEmailGenerator(llmProvider),
		llm.NewCandidateComparator(llmProvider),
		llm.NewJobParser(llmProvider),
		candidateRepo,
		jobRepo,
		commRepo,
		hiringGRPC,
		auditor,
	)

	temporalWorker := temporal.NewWorker(zapLog.Log, temporalClient, &conf.Temporal, activities)

	mqClient := mq.NewRabbitMQ(zapLog.Log, &conf.RabbitMQ)

	candidateHandler := worker.NewCandidateConsumerHandler(zapLog.Log, temporalClient)
	mqConsumer := mq.NewMQConsumer(zapLog.Log, mqClient.Conn, candidateHandler,
		mq.WithQueueName("hiring.candidate.created"),
		mq.WithConsumerExchange("hiring.events"),
		mq.WithConsumerRoutingKey("candidate.created"),
		mq.WithQuorumQueue(),
		mq.WithDLX("hiring.dlx", "hiring.candidate.dead"),
		mq.WithMessageTTL(300_000),
		mq.WithMaxDeliveries(3),
		mq.WithPrefetchCount(1),
	)

	dlqHandler := worker.NewDLQHandler(zapLog.Log, temporalClient)
	dlqWorker := mq.NewMQConsumer(zapLog.Log, mqClient.Conn, dlqHandler,
		mq.WithQueueName("hiring.candidate.dead"),
		mq.WithConsumerExchange("hiring.dlx"),
		mq.WithConsumerRoutingKey("hiring.candidate.dead"),
	)

	return &infrastructureComponents{
		cfg:            conf,
		log:            zapLog,
		pool:           pool,
		cloudinary:     cloudinary,
		llmClient:      llmClient,
		temporalClient: temporalClient,
		temporalWorker: temporalWorker,
		mqClient:       mqClient,
		mqConsumer:     mqConsumer,
		dlqWorker:      dlqWorker,
		hiringGRPC:     hiringGRPC,
	}, nil
}
