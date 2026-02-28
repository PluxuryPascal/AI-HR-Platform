package main

import (
	"backend/internal/cache"
	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repo"
	"backend/internal/server"
	"backend/internal/server/router/invite"
	"backend/internal/server/router/user"
	"backend/internal/usecase"
	"backend/pkg/config"
	"backend/pkg/hash"
	"backend/pkg/logger"
	"backend/pkg/svc"
	"backend/pkg/token"
	"context"
	"fmt"
	"log"

	"github.com/lestrrat-go/jwx/v2/jwk"
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
	auth   *handler.AuthHandler
	invite *handler.InviteHandler
}

type infrastructureComponents struct {
	cfg       *config.Config
	log       *logger.Log
	pool      *db.PostgresClient
	redisPool *db.RedisClient
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
	usecases := initUseCases(infra, utils, repos)
	handlers, sessionMiddleware := initHandlers(infra, utils, usecases)

	apiServer := createApiServer(ctx, infra.cfg, infra.log, handlers, sessionMiddleware)

	if err := svc.Run(ctx, infra.log.Log, []svc.Service{
		infra.log,
		infra.pool,
		infra.redisPool,
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

	return &infrastructureComponents{
		cfg:       conf,
		log:       zapLog,
		pool:      pool,
		redisPool: redisPool,
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

func initUseCases(infra *infrastructureComponents, utils *utilityComponents, r repos) usecases {
	return usecases{
		auth:   usecase.NewAuthUseCase(r.user, utils.cacheManager, utils.t, utils.h),
		invite: usecase.NewInviteUseCase(infra.cfg, r.invite, utils.cacheManager, utils.t, utils.h),
	}
}

func initHandlers(infra *infrastructureComponents, utils *utilityComponents, u usecases) (handlers, middleware.SessionMiddleware) {
	sessionMiddleware := middleware.NewSessionMiddleware(
		infra.log,
		infra.redisPool,
		utils.cacheManager,
	)

	h := handlers{
		auth:   handler.NewAuthHandler(&infra.cfg.Server, infra.log.Log, u.auth),
		invite: handler.NewInviteHandler(&infra.cfg.Server, infra.log.Log, u.invite),
	}

	return h, sessionMiddleware
}

func createApiServer(ctx context.Context, cfg *config.Config, log *logger.Log, h handlers, sessionMiddleware middleware.SessionMiddleware) *server.Api {
	return server.NewApiServer(
		&cfg.Server,
		server.WithLogger(log.Log),
		server.WithRouterGroup(ctx, "/auth",
			user.NewRouter(
				h.auth,
				sessionMiddleware.RateLimit(cfg.RateLimit["auth"]),
			),
		),
		server.WithRouterGroup(ctx, "/invite",
			invite.NewRouter(
				h.invite,
				sessionMiddleware.RateLimit(cfg.RateLimit["invite"]),
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
