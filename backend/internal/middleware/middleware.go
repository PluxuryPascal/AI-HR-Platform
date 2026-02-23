package middleware

import (
	"backend/internal/cache"
	"backend/internal/db"
	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/token"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type SessionMiddleware interface {
	RateLimit(rateLimit config.RateLimit) echo.MiddlewareFunc
	Session(t *token.JWTtoken) echo.MiddlewareFunc
}

type middleware struct {
	log          *zap.Logger
	redisClient  *db.RedisClient
	cacheManager *cache.Manager
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
			c.Set("group", subject.Group)

			return next(c)
		}
	}
}

func NewSessionMiddleware(log *logger.Log, redisClient *db.RedisClient, cacheManager *cache.Manager) SessionMiddleware {
	return &middleware{
		log:          log.Log,
		redisClient:  redisClient,
		cacheManager: cacheManager,
	}
}

var _ SessionMiddleware = (*middleware)(nil)
