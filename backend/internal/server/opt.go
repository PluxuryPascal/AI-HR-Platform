package server

import (
	"backend/pkg/router"
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// WithRouter регистрирует роутеры в базовой группе /api.
func WithRouter(ctx context.Context, routes ...router.Router) func(*Api) {
	return func(s *Api) {
		g := s.api.Group("/api/v1")
		for _, apiRouter := range routes {
			for _, route := range apiRouter.Routes() {
				route.Register(ctx, g)
			}
		}
	}
}

// WithRouterGroup регистрирует роутеры в группе /api/{prefix},
// например WithRouterGroup(ctx, "/v1", userRouter, inviteRouter).
func WithRouterGroup(ctx context.Context, prefix string, routes ...router.Router) func(*Api) {
	return func(s *Api) {
		g := s.api.Group("/api/v1" + prefix)
		for _, apiRouter := range routes {
			for _, route := range apiRouter.Routes() {
				route.Register(ctx, g)
			}
		}
	}
}

func WithLogger(log *zap.Logger) func(*Api) {
	return func(s *Api) {
		s.api.Use(middleware.RequestLoggerWithConfig(
			middleware.RequestLoggerConfig{
				LogURI:       true,
				LogStatus:    true,
				LogError:     true,
				LogLatency:   true,
				LogRemoteIP:  true,
				LogMethod:    true,
				LogRequestID: true,

				LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
					fields := []zap.Field{
						zap.String("method", v.Method),
						zap.String("uri", v.URI),
						zap.Int("status", v.Status),
						zap.Duration("latency", v.Latency),
						zap.String("remote_ip", v.RemoteIP),
					}

					if v.Error == nil {
						log.Info("request", fields...)

						return nil
					}

					var httpErr *echo.HTTPError
					if errors.As(v.Error, &httpErr) {
						fields = append(fields, zap.Any("ErrorMessage", httpErr.Message))
						log.Error("request", fields...)
					} else {
						fields = append(fields, zap.Error(v.Error))
						log.Error("request", fields...)
					}

					return nil
				},
			},
		))
	}
}

func WithMiddleware(middlewares ...echo.MiddlewareFunc) func(*Api) {
	return func(s *Api) {
		s.api.Use(middlewares...)
	}
}
