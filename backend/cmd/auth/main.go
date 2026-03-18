package main

import (
	"backend/internal/audit"
	"backend/internal/cache"
	"backend/internal/db"
	auth_grpc "backend/internal/grpc-handler/auth"
	httpHandler "backend/internal/http-handler"
	"backend/internal/middleware"
	authv1 "backend/internal/proto/auth/v1"
	pb "backend/internal/proto/hiring/v1"
	"backend/internal/repo"
	"backend/internal/server"
	"backend/internal/server/router/invite"
	"backend/internal/server/router/user"
	"backend/internal/usecase"
	"backend/internal/worker"
	"backend/pkg/config"
	"backend/pkg/grpc"
	"backend/pkg/hash"
	"backend/pkg/logger"
	"backend/pkg/rbac"
	"backend/pkg/svc"
	"backend/pkg/token"
	"context"
	"fmt"
	"log"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"go.uber.org/zap"
)

type repos struct {
	user   repo.UserRepository
	invite repo.InviteRepository
}

type usecases struct {
	auth   usecase.AuthUseCase
	invite usecase.InviteUseCase
}

type handlers struct {
	auth   *httpHandler.AuthHandler
	invite *httpHandler.InviteHandler
}

type infrastructureComponents struct {
	cfg       *config.Config
	log       *logger.Log
	pool      *db.PostgresClient
	redisPool *db.RedisClient
	casbin    *rbac.CasbinClient
	grpcCl    *grpc.Client
	grpcSrv   *grpc.Server
}

type utilityComponents struct {
	cacheManager *cache.Manager
	t            *token.JWTtoken
	h            *hash.Argon2
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

	var hiringClient pb.HiringServiceClient
	infra.grpcCl.OnInit(func(c *grpc.Client) {
		hiringClient = pb.NewHiringServiceClient(c.GetConn())
	})

	auditor := audit.NewLogger(infra.log.Log, infra.pool)
	if err := auditor.SeedActionTypes(ctx); err != nil {
		infra.log.Log.Warn("failed to seed audit action types", zap.Error(err))
	}

	usecases := initUseCases(infra, utils, repos, &hiringClient, auditor)
	initGrpcHandlers(infra, repos)
	handlers, sessionMiddleware := initHandlers(infra, utils, usecases)

	apiServer := createApiServer(ctx, infra.cfg, utils.t, infra.log, handlers, sessionMiddleware)
	recoveryWorker := worker.NewInviteRecoveryWorker(infra.log.Log, &infra.cfg.InviteRecovery, usecases.invite)

	if err := svc.Run(ctx, infra.log.Log, []svc.Service{
		infra.log,
		infra.pool,
		infra.redisPool,
		infra.casbin,
		apiServer,
		infra.grpcCl,
		infra.grpcSrv,
		recoveryWorker,
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

	hiringClientCfg := conf.GRPC.Clients["hiring"]
	grpcClient := grpc.NewClient("hiring", zapLog.Log, &hiringClientCfg)

	serverCfg := conf.GRPC.Servers["auth"]
	grpcServer := grpc.NewServer("auth", zapLog.Log, &serverCfg)

	return &infrastructureComponents{
		cfg:       conf,
		log:       zapLog,
		pool:      pool,
		redisPool: redisPool,
		casbin:    casbinClient,
		grpcCl:    grpcClient,
		grpcSrv:   grpcServer,
	}, nil
}

func initUtilities(infra *infrastructureComponents) (*utilityComponents, error) {
	prvKey, err := loadKey(infra.cfg.Token.PrivateKey.Path)
	if err != nil {
		return nil, fmt.Errorf("load key error: %w", err)
	}

	t, err := token.NewJWTtoken(infra.cfg.Token.Issuer, infra.cfg.Token.ExpireAt, prvKey)
	if err != nil {
		return nil, fmt.Errorf("create token error: %w", err)
	}

	h := hash.NewArgon2(infra.cfg.Hash)
	cacheManager := cache.NewManager(infra.redisPool, cache.WithPrefix("ai_hr"))

	return &utilityComponents{
		cacheManager: cacheManager,
		t:            t,
		h:            h,
	}, nil
}

func initRepositories(infra *infrastructureComponents) repos {
	return repos{
		user:   repo.NewUserRepo(infra.pool),
		invite: repo.NewInviteRepo(infra.pool),
	}
}

func initUseCases(infra *infrastructureComponents, utils *utilityComponents, r repos, hiringClient *pb.HiringServiceClient, auditor *audit.Logger) usecases {
	return usecases{
		auth:   usecase.NewAuthUseCase(r.user, utils.cacheManager, utils.t, utils.h, infra.casbin, auditor),
		invite: usecase.NewInviteUseCase(infra.cfg, r.invite, r.user, utils.cacheManager, utils.t, utils.h, infra.casbin, hiringClient, auditor),
	}
}

func initGrpcHandlers(infra *infrastructureComponents, r repos) {
	authHandler := auth_grpc.NewAuthHandler(infra.log.Log, r.user)
	infra.grpcSrv.OnInit(func(s *grpc.Server) {
		authv1.RegisterAuthServiceServer(s.GetServer(), authHandler)
	})
}

func initHandlers(infra *infrastructureComponents, utils *utilityComponents, u usecases) (handlers, middleware.Middleware) {
	middleware := middleware.NewMiddleware(
		infra.log,
		infra.redisPool,
		utils.cacheManager,
		infra.casbin,
	)

	h := handlers{
		auth:   httpHandler.NewAuthHandler(&infra.cfg.Server, infra.log.Log, u.auth),
		invite: httpHandler.NewInviteHandler(&infra.cfg.Server, infra.log.Log, u.invite),
	}

	return h, middleware
}

func createApiServer(ctx context.Context, cfg *config.Config, t *token.JWTtoken, log *logger.Log, h handlers, middleware middleware.Middleware) *server.Api {
	return server.NewApiServer(
		cfg.Server.Ports["auth"],
		&cfg.Server,
		server.WithLogger(log.Log),
		server.WithRouterGroup(ctx, "/auth",
			user.NewRouter(
				h.auth,
				middleware.RateLimit("auth", cfg.RateLimit["auth"]),
				middleware.Session(t),
			),
		),
		server.WithRouterGroup(ctx, "/invite",
			invite.NewRouter(
				h.invite,
				middleware.RateLimit("invite", cfg.RateLimit["invite"]),
				middleware.Session(t),
				middleware.RBAC(),
			),
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
