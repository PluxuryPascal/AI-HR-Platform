package main

import (
	"backend/internal/audit"
	"backend/internal/cache"
	"backend/internal/db"
	grpchiring "backend/internal/grpc-handler/hiring"
	handler "backend/internal/http-handler"
	"backend/internal/middleware"
	pbAI "backend/internal/proto/ai_engine/v1"
	pbAuth "backend/internal/proto/auth/v1"
	pb "backend/internal/proto/hiring/v1"
	"backend/internal/repo"
	"backend/internal/server"
	"backend/internal/temporal"
	"backend/internal/server/router/access"
	"backend/internal/server/router/ai_settings"
	"backend/internal/server/router/candidate"
	"backend/internal/server/router/chat"
	"backend/internal/server/router/dashboard"
	"backend/internal/server/router/department"
	"backend/internal/server/router/interview"
	"backend/internal/server/router/job"
	"backend/internal/server/router/pipeline"
	"backend/internal/openrouter"
	"backend/internal/usecase"
	"backend/pkg/config"
	"backend/pkg/grpc"
	"backend/pkg/logger"
	"backend/pkg/mq"
	"backend/pkg/rbac"
	"backend/pkg/storage"
	"backend/pkg/svc"
	"backend/pkg/token"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"go.uber.org/zap"
)

type repos struct {
	access     repo.AccessRepository
	department repo.DepartmentRepository
	job        repo.JobRepository
	candidate  repo.CandidateRepository
	pipeline   repo.PipelineRepository
	dashboard  repo.DashboardRepository
	aiSettings repo.AiSettingsRepository
	chat       repo.ChatRepository
}

type usecases struct {
	access     usecase.AccessUseCase
	department usecase.DepartmentUseCase
	job        usecase.JobUseCase
	candidate  usecase.CandidateUseCase
	pipeline   usecase.PipelineUseCase
	dashboard  usecase.DashboardUseCase
	chat       usecase.ChatUseCase
	interview  usecase.InterviewUseCase
}

type handlers struct {
	job        *handler.JobHandler
	candidate  *handler.CandidateHandler
	access     *handler.AccessHandler
	pipeline   *handler.PipelineHandler
	department *handler.DepartmentHandler
	dashboard  *handler.DashboardHandler
	aiSettings *handler.AiSettingsHandler
	chat       *handler.ChatHandler
	interview  *handler.InterviewHandler
}

type infrastructureComponents struct {
	cfg          *config.Config
	log          *logger.Log
	pool         *db.PostgresClient
	redisPool    *db.RedisClient
	grpcServer   *grpc.Server
	authClient   *grpc.Client
	aiClient     *grpc.Client
	storage      *storage.CloudinaryStorage
	casbinClient *rbac.CasbinClient
	mqClient     *mq.RabbitMQ
	mqPublisher  *mq.MQPublisher
	openRouter   *openrouter.Client
	temporalClient *temporal.Client
}

type utilityComponents struct {
	cacheManager *cache.Manager
	token        *token.JWTtoken
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

	auditor := audit.NewLogger(infra.log.Log, infra.pool)
	infra.pool.AddAfterRun(func(ctx context.Context, _ *pgxpool.Pool) error {
		if err := auditor.SeedActionTypes(ctx); err != nil {
			infra.log.Log.Warn("failed to seed audit action types", zap.Error(err))
		}
		return nil
	})

	usecases := initUseCases(infra, utils, repos, auditor)

	handlers, sessionMiddleware := initHandlers(infra, utils, usecases, repos, auditor)

	grpcHandler := grpchiring.NewHandler(infra.log.Log, usecases.access, usecases.candidate)

	infra.grpcServer.OnInit(func(s *grpc.Server) {
		pb.RegisterHiringServiceServer(s.GetServer(), grpcHandler)
	})

	apiServer := createApiServer(ctx, infra.cfg, infra.log, handlers, sessionMiddleware, utils.token)

	if err := svc.Run(ctx, infra.log.Log, []svc.Service{
		infra.log,
		infra.pool,
		infra.redisPool,
		infra.grpcServer,
		infra.mqClient,
		infra.mqPublisher,
		infra.storage,
		infra.casbinClient,
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

	casbinClient := rbac.NewCasbinClient(zapLog.Log, pool.ConnConfig(), "casbin/model.conf")

	storage := storage.NewCloudinaryStorage(zapLog.Log, &conf.Cloudinary)

	mqClient := mq.NewRabbitMQ(zapLog.Log, &conf.RabbitMQ)
	mqPublisher := mq.NewMQPublisher(zapLog.Log, mqClient)

	serverCfg := conf.GRPC.Servers["hiring"]
	grpcServer := grpc.NewServer("hiring", zapLog.Log, &serverCfg)

	authClientCfg := conf.GRPC.Clients["auth"]
	authGrpcClient := grpc.NewClient("auth", zapLog.Log, &authClientCfg)

	aiClientCfg := conf.GRPC.Clients["ai_engine"]
	aiGrpcClient := grpc.NewClient("ai_engine", zapLog.Log, &aiClientCfg)

	openRouter := openrouter.NewClient(zapLog.Log, conf.OpenRouter.BaseURL)

	return &infrastructureComponents{
		cfg:          conf,
		log:          zapLog,
		pool:         pool,
		redisPool:    redisPool,
		grpcServer:   grpcServer,
		authClient:   authGrpcClient,
		aiClient:     aiGrpcClient,
		storage:      storage,
		casbinClient: casbinClient,
		mqClient:     mqClient,
		mqPublisher:  mqPublisher,
		openRouter:   openRouter,
		temporalClient: temporal.NewClient(zapLog.Log, nil, &conf.Temporal),
	}, nil
}

func initUtilities(infra *infrastructureComponents) (*utilityComponents, error) {
	cacheManager := cache.NewManager(infra.redisPool, cache.WithPrefix("ai_hr"))

	prvKey, err := loadKey(infra.cfg.Token.PrivateKey.Path)
	if err != nil {
		return nil, fmt.Errorf("load key error: %w", err)
	}

	t, err := token.NewJWTtoken(infra.cfg.Token.Issuer, infra.cfg.Token.ExpireAt, prvKey)
	if err != nil {
		return nil, fmt.Errorf("create jwt token error: %w", err)
	}

	return &utilityComponents{
		cacheManager: cacheManager,
		token:        t,
	}, nil
}

func initRepositories(infra *infrastructureComponents) repos {
	return repos{
		access:     repo.NewAccessRepository(infra.pool),
		department: repo.NewDepartmentRepo(infra.pool),
		job:        repo.NewJobRepo(infra.pool),
		candidate:  repo.NewCandidateRepo(infra.pool),
		pipeline:   repo.NewPipelineRepo(infra.pool),
		dashboard:  repo.NewDashboardRepo(infra.pool),
		aiSettings: repo.NewAiSettingsRepo(infra.pool),
		chat:       repo.NewChatRepo(infra.pool),
	}
}

func initUseCases(infra *infrastructureComponents, utils *utilityComponents, r repos, auditor *audit.Logger) usecases {
	var authSvcClient pbAuth.AuthServiceClient
	infra.authClient.OnInit(func(c *grpc.Client) {
		authSvcClient = pbAuth.NewAuthServiceClient(c.GetConn())
	})

	var aiSvcClient pbAI.AIEngineServiceClient
	infra.aiClient.OnInit(func(c *grpc.Client) {
		aiSvcClient = pbAI.NewAIEngineServiceClient(c.GetConn())
	})

	return usecases{
		access:     usecase.NewAccessUseCase(r.access, auditor),
		department: usecase.NewDepartmentUseCase(r.department),
		job:        usecase.NewJobUseCase(r.job, r.access, auditor),
		candidate:  usecase.NewCandidateUseCase(infra.log.Log, r.candidate, r.pipeline, r.job, r.access, infra.storage, infra.mqPublisher, auditor),
		pipeline:   usecase.NewPipelineUseCase(r.pipeline, auditor),
		dashboard:  usecase.NewDashboardUseCase(infra.log.Log, r.dashboard, authSvcClient, aiSvcClient),
		chat:       usecase.NewChatUseCase(infra.log.Log, r.chat, infra.temporalClient),
		interview:  usecase.NewInterviewUseCase(infra.log.Log, infra.temporalClient.TemporalClient),
	}
}

func initHandlers(infra *infrastructureComponents, utils *utilityComponents, u usecases, r repos, auditor *audit.Logger) (handlers, middleware.Middleware) {
	mw := middleware.NewMiddleware(
		infra.log,
		infra.redisPool,
		utils.cacheManager,
		infra.casbinClient,
	)

	h := handlers{
		job:        handler.NewJobHandler(infra.log.Log, u.job),
		candidate:  handler.NewCandidateHandler(infra.log.Log, u.candidate),
		access:     handler.NewAccessHandler(infra.log.Log, u.access),
		pipeline:   handler.NewPipelineHandler(infra.log.Log, u.pipeline),
		department: handler.NewDepartmentHandler(infra.log.Log, u.department),
		dashboard:  handler.NewDashboardHandler(infra.log.Log, u.dashboard),
		aiSettings: handler.NewAiSettingsHandler(infra.log.Log, r.aiSettings, infra.openRouter),
		chat:       handler.NewChatHandler(infra.log.Log, u.chat),
	}

	return h, mw
}

func createApiServer(ctx context.Context, cfg *config.Config, log *logger.Log, h handlers, mw middleware.Middleware, t *token.JWTtoken) *server.Api {
	return server.NewApiServer(
		cfg.Server.Ports["hiring"],
		&cfg.Server,
		server.WithLogger(log.Log),
		server.WithRouterGroup(ctx, "/jobs",
			job.NewRouter(h.job, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/jobs/:job_id/candidates",
			candidate.NewJobScopedRouter(h.candidate, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/jobs/:job_id/stages",
			pipeline.NewRouter(h.pipeline, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/candidates",
			candidate.NewRouter(h.candidate, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/jobs/:job_id/access",
			access.NewRouter(h.access, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/departments",
			department.NewRouter(h.department, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/dashboard",
			dashboard.NewRouter(h.dashboard, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/ai-settings",
			ai_settings.NewRouter(h.aiSettings, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/chat",
			chat.NewRouter(h.chat, mw.Session(t), mw.RBAC()),
		),
		server.WithRouterGroup(ctx, "/interview",
			interview.NewRouter(h.interview, mw.Session(t), mw.RBAC()),
		),
	)
}

func loadKey(path string) (jwk.Key, error) {
	keySet, err := jwk.ReadFile(path, jwk.WithPEM(true))
	if err != nil {
		return nil, fmt.Errorf("read key error: %w", err)
	}

	key, ok := keySet.Key(0)
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return key, nil
}
