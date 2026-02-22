package main

import (
	"backend/internal/cache"
	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/repo"
	"backend/internal/server"
	"backend/internal/server/router/user"
	"backend/pkg/config"
	"backend/pkg/hash"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"backend/pkg/svc"
	"backend/pkg/token"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type repos struct {
	user repo.UserRepository
}

type handlers struct {
	auth *handler.AuthHandler
}

type coreComponents struct {
	cfg          *config.Config
	log          *logger.Log
	pool         *db.Client
	redisPool    *db.RedisClient
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
	core, err := initCoreComponents()
	if err != nil {
		return fmt.Errorf("init core components error: %w", err)
	}

	handlers, middleware := initServices(core)

	apiServer := createApiServer(ctx, core, handlers, middleware)

	if err := svc.Run(ctx, []svc.Service{
		core.log,
		core.pool,
		core.redisPool,
		apiServer,
	}); err != nil {
		return fmt.Errorf("run service error: %w", err)
	}

	return nil
}

func initCoreComponents() (*coreComponents, error) {
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

	prvKey, err := loadKey(conf.Token.PrivateKey.Path)
	if err != nil {
		return nil, fmt.Errorf("load key error: %w", err)
	}

	t := token.NewJWTtoken(conf.Token.Issuer, conf.Token.ExpireAt, prvKey)

	h := hash.NewArgon2(conf.Hash)

	cacheManager := cache.NewManager(redisPool, cache.WithPrefix("ai_hr"))

	return &coreComponents{
		cfg:          conf,
		log:          zapLog,
		pool:         pool,
		redisPool:    redisPool,
		cacheManager: cacheManager,
		t:            t,
		h:            h,
	}, nil
}

func initServices(c *coreComponents) (handlers, middleware.SessionMiddleware) {
	sessionMiddleware := middleware.NewSessionMiddleware(
		c.redisPool,
		c.cacheManager,
	)

	repos := repos{
		user: repo.NewUserRepo(c.pool),
	}

	handlers := handlers{
		auth: handler.NewAuthHandler(c.log.Log, repos.user, c.cacheManager, c.t, c.h),
	}

	return handlers, sessionMiddleware
}

func createApiServer(ctx context.Context, core *coreComponents, handlers handlers, middleware middleware.SessionMiddleware) *server.Api {
	return server.NewApiServer(
		core.cfg.Server.Port,
		server.WithLogger(core.log.Log),
		server.WithRouter(ctx,
			user.NewRouter(handlers.auth, middleware.RateLimit(5, 5*time.Minute)),
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
