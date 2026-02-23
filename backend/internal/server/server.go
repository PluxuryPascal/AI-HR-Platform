package server

import (
	"backend/internal/server/validator"
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Api struct {
	log  *zap.Logger
	api  *echo.Echo
	port int
}

func (a *Api) DependsOn() []string {
	return []string{"logger", "db", "redis"}
}

func (a *Api) HealthCheck(ctx context.Context) error {
	return nil
}

func (a *Api) Init(ctx context.Context) error {
	return nil
}

func (a *Api) Name() string {
	return "api"
}

func (a *Api) Run(ctx context.Context) error {
	if err := a.api.Start(fmt.Sprintf(":%d", a.port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start API server on port %d: %w", a.port, err)
	}

	return nil
}

func (a *Api) Stop(ctx context.Context) error {
	if err := a.api.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown API server gracefully: %w", err)
	}

	return nil
}

func NewApiServer(cfg *config.Server, opts ...func(*Api)) *Api {
	e := echo.New()

	e.Validator = validator.NewValidator()
	e.Use(middleware.Recover())
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout

	api := &Api{
		api:  e,
		port: cfg.Port,
	}

	for _, opt := range opts {
		opt(api)
	}

	return api
}

var _ svc.Service = (*Api)(nil)
