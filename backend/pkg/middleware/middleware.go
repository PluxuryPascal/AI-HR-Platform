package middleware

import (
	"backend/internal/cache"
	"backend/internal/db"
	"backend/pkg/token"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type sessionSubject struct {
	UserID string `json:"user_id"`
	Group  string `json:"group"`
}

type SessionMiddleware interface {
	RateLimit(requests int, window time.Duration) echo.MiddlewareFunc
	Session(t *token.JWTtoken) echo.MiddlewareFunc
}

type middleware struct {
	redisClient  *db.RedisClient
	cacheManager *cache.Manager
}

func (m *middleware) RateLimit(requests int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ipKey := c.RealIP()

			count, err := cache.IncrWithTTL(c.Request().Context(), m.cacheManager, cache.RateLimitKey, ipKey, window)
			if err != nil {
				return next(c)
			}

			if count == 1 {
				c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
			}

			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(requests))

			remaining := requests - int(count)
			if remaining < 0 {
				remaining = 0
			}

			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if count > int64(requests) {
				return echo.NewHTTPError(http.StatusTooManyRequests, map[string]string{
					"message": "too many requests",
					"reset":   strconv.FormatInt(time.Now().Add(window).Unix(), 10),
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
				return echo.NewHTTPError(http.StatusUnauthorized, "cookie not found: %w", err)
			}

			pubKey, err := t.PrvKey.PublicKey()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "public key error: %w", err)
			}

			token, err := t.VerifyToken([]byte(cookie.Value), pubKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "token not valid: %w", err)
			}

			rawSubject, err := cache.Get(c.Request().Context(), m.cacheManager, cache.SessionKey, token.Subject())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "session not found: %w", err)
			}

			var subject sessionSubject
			if err := json.Unmarshal([]byte(rawSubject), &subject); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "unmarshal subject error: %w", err)
			}

			c.Set("id", subject.UserID)
			c.Set("group", subject.Group)

			return next(c)
		}
	}
}

func NewSessionMiddleware(redisClient *db.RedisClient, cacheManager *cache.Manager) SessionMiddleware {
	return &middleware{
		redisClient:  redisClient,
		cacheManager: cacheManager,
	}
}

var _ SessionMiddleware = (*middleware)(nil)
