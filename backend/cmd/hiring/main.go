package main

import (
	"backend/internal/cache"
	"backend/internal/db"
	grpchiring "backend/internal/grpc-handler/hiring"
	"backend/internal/middleware"
	pb "backend/internal/proto/hiring/v1"
	"backend/internal/repo"
	"backend/internal/server"
	"backend/internal/usecase"
	"backend/pkg/config"
	"backend/pkg/grpc"
	"backend/pkg/logger"
	"backend/pkg/rbac"
	"backend/pkg/svc"
	"context"
	"fmt"
	"log"
)

type repos struct {
	access     repo.AccessRepository
	department repo.DepartmentRepository
	job        repo.JobRepository
	candidate  repo.CandidateRepository
}

type usecases struct {
	access     usecase.AccessUseCase
	department usecase.DepartmentUseCase
	job        usecase.JobUseCase
	candidate  usecase.CandidateUseCase
}

type handlers struct {
}

type infrastructureComponents struct {
	cfg        *config.Config
	log        *logger.Log
	pool       *db.PostgresClient
	redisPool  *db.RedisClient
	grpcServer *grpc.Server
}

type utilityComponents struct {
	casbinClient *rbac.CasbinClient
	cacheManager *cache.Manager
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

	utils, err := initUtilities(infra)
	if err != nil {
		return fmt.Errorf("init utilities error: %w", err)
	}

	repos := initRepositories(infra)
	usecases := initUseCases(infra, utils, repos)
	handlers, sessionMiddleware := initHandlers(infra, utils, usecases)

	grpcHandler := grpchiring.NewHandler(infra.log.Log, usecases.access)

	infra.grpcServer.OnInit(func(s *grpc.Server) {
		pb.RegisterHiringServiceServer(s.GetServer(), grpcHandler)
	})

	apiServer := createApiServer(ctx, infra.cfg, infra.log, handlers, sessionMiddleware)

	if err := svc.Run(ctx, infra.log.Log, []svc.Service{
		infra.log,
		infra.pool,
		infra.redisPool,
		infra.grpcServer,
		apiServer,
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

	serverCfg := conf.GRPC.Servers["hiring"]
	grpcServer := grpc.NewServer("hiring", zapLog.Log, &serverCfg)

	return &infrastructureComponents{
		cfg:        conf,
		log:        zapLog,
		pool:       pool,
		redisPool:  redisPool,
		grpcServer: grpcServer,
	}, nil
}

func initUtilities(infra *infrastructureComponents) (*utilityComponents, error) {
	casbinClient := rbac.NewCasbinClient(infra.log.Log, infra.pool.ConnConfig(), "casbin/model.conf")

	cacheManager := cache.NewManager(infra.redisPool, cache.WithPrefix("ai_hr"))

	return &utilityComponents{
		casbinClient: casbinClient,
		cacheManager: cacheManager,
	}, nil
}

func initRepositories(infra *infrastructureComponents) repos {
	return repos{
		access:     repo.NewAccessRepository(infra.pool),
		department: repo.NewDepartmentRepo(infra.pool),
		job:        repo.NewJobRepo(infra.pool),
		candidate:  repo.NewCandidateRepo(infra.pool),
	}
}

func initUseCases(infra *infrastructureComponents, utils *utilityComponents, r repos) usecases {
	return usecases{
		access:     usecase.NewAccessUseCase(r.access),
		department: usecase.NewDepartmentUseCase(r.department),
		job:        usecase.NewJobUseCase(r.job),
		candidate:  usecase.NewCandidateUseCase(r.candidate),
	}
}

func initHandlers(infra *infrastructureComponents, utils *utilityComponents, u usecases) (handlers, middleware.Middleware) {
	middleware := middleware.NewMiddleware(
		infra.log,
		infra.redisPool,
		utils.cacheManager,
		utils.casbinClient,
	)

	h := handlers{}

	return h, middleware
}

func createApiServer(ctx context.Context, cfg *config.Config, log *logger.Log, h handlers, middleware middleware.Middleware) *server.Api {
	return server.NewApiServer(
		cfg.Server.Ports["hiring"],
		&cfg.Server,
		server.WithLogger(log.Log),
	)
}
