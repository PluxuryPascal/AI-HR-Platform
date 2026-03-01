package middleware

import (
	"backend/internal/cache"
	"backend/internal/db"
	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/rbac"
	"backend/pkg/token"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Middleware interface {
	RateLimit(rateLimit config.RateLimit) echo.MiddlewareFunc
	Session(t *token.JWTtoken) echo.MiddlewareFunc
	RBAC() echo.MiddlewareFunc
}

type middleware struct {
	log            *zap.Logger
	redisClient    *db.RedisClient
	cacheManager   *cache.Manager
	casbinEnforcer *rbac.CasbinClient
}

func (m *middleware) RateLimit(rateLimit config.RateLimit) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ipKey := c.RealIP()

			count, err := cache.IncrWithTTL(c.Request().Context(), m.cacheManager, cache.RateLimitKey, ipKey, rateLimit.Window)
			if err != nil {
				m.log.Error("rate limit error", zap.Error(err))
				return echo.NewHTTPError(http.StatusServiceUnavailable, "service temporarily unavailable")
			}

			if count == 1 {
				c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rateLimit.Window).Unix(), 10))
			}

			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(rateLimit.Requests))

			remaining := rateLimit.Requests - int(count)
			if remaining < 0 {
				remaining = 0
			}

			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if count > int64(rateLimit.Requests) {
				return echo.NewHTTPError(http.StatusTooManyRequests, map[string]string{
					"message": "too many requests",
					"reset":   strconv.FormatInt(time.Now().Add(rateLimit.Window).Unix(), 10),
				})
			}

			return next(c)
		}
	}
}

func (m *middleware) Session(t *token.JWTtoken) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("access_token")
			if err != nil {
				m.log.Warn("cookie not found", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "cookie not found")
			}

			token, err := t.VerifyToken([]byte(cookie.Value))
			if err != nil {
				m.log.Warn("token not valid", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "token not valid")
			}

			subject, err := cache.Get(c.Request().Context(), m.cacheManager, cache.SessionKey, token.Subject())
			if err != nil {
				m.log.Warn("session not found", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "session not found")
			}

			c.Set("id", subject.UserID)
			c.Set("team_id", subject.TeamID)
			c.Set("role", subject.Role)

			return next(c)
		}
	}
}

func (m *middleware) RBAC() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("id").(string)
			teamID := c.Get("team_id").(string)

			obj := c.Path()
			act := c.Request().Method

			ok, err := m.casbinEnforcer.Enforce(userID, teamID, obj, act)
			if err != nil {
				m.log.Error("rbac error", zap.Error(err))
				return echo.NewHTTPError(http.StatusServiceUnavailable, "service temporarily unavailable")
			}

			if !ok {
				m.log.Warn("rbac denied", zap.String("user_id", userID), zap.String("team_id", teamID), zap.String("obj", obj), zap.String("act", act))
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			return next(c)
		}
	}
}

func NewMiddleware(log *logger.Log, redisClient *db.RedisClient, cacheManager *cache.Manager, casbinEnforcer *rbac.CasbinClient) Middleware {
	return &middleware{
		log:            log.Log,
		redisClient:    redisClient,
		cacheManager:   cacheManager,
		casbinEnforcer: casbinEnforcer,
	}
}

var _ Middleware = (*middleware)(nil)
