package main

import (
	"backend/internal/db"
	grpcaiengine "backend/internal/grpc-handler/ai_engine"
	pb "backend/internal/proto/ai_engine/v1"
	"backend/internal/repo"
	"backend/internal/temporal"
	"backend/pkg/config"
	"backend/pkg/grpc"
	"backend/pkg/logger"
	"backend/pkg/svc"
	"context"
	"fmt"
	"log"
)

type infrastructureComponents struct {
	cfg            *config.Config
	log            *logger.Log
	pool           *db.PostgresClient
	redisPool      *db.RedisClient
	temporalClient *temporal.Client
	grpcServer     *grpc.Server
}

type repos struct {
	candidate     repo.CandidateRepository
	communication repo.CommunicationRepository
}

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("application error: %v", err)
	}

	log.Println("success shutdown")
}

func run(ctx context.Context) error {
	infra, err := initInfrastructure()
	if err != nil {
		return fmt.Errorf("init infrastructure error: %w", err)
	}

	repos := initRepositories(infra)

	grpcHandler := grpcaiengine.NewHandler(
		infra.log.Log,
		infra.temporalClient,
		repos.candidate,
		repos.communication,
	)

	infra.grpcServer.OnInit(func(s *grpc.Server) {
		pb.RegisterAIEngineServiceServer(s.GetServer(), grpcHandler)
	})

	if err := svc.Run(ctx, infra.log.Log, []svc.Service{
		infra.log,
		infra.pool,
		infra.redisPool,
		infra.temporalClient,
		infra.grpcServer,
	}); err != nil {
		return fmt.Errorf("run service error: %w", err)
	}

	return nil
}

func initInfrastructure() (*infrastructureComponents, error) {
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

	redisPool, err := db.NewRedis(zapLog.Log, conf)
	if err != nil {
		return nil, fmt.Errorf("create redis error: %w", err)
	}

	temporalClient := temporal.NewClient(zapLog.Log, nil, &conf.Temporal)

	serverCfg := conf.GRPC.Servers["ai_engine"]
	grpcServer := grpc.NewServer("ai_engine", zapLog.Log, &serverCfg)

	return &infrastructureComponents{
		cfg:            conf,
		log:            zapLog,
		pool:           pool,
		redisPool:      redisPool,
		temporalClient: temporalClient,
		grpcServer:     grpcServer,
	}, nil
}

func initRepositories(infra *infrastructureComponents) repos {
	return repos{
		candidate:     repo.NewCandidateRepo(infra.pool),
		communication: repo.NewCommunicationRepo(infra.pool),
	}
}
