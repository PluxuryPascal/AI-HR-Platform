package handler

import (
	"backend/internal/cache"
	"backend/internal/domain"
	"backend/internal/repo"
	"backend/pkg/hash"
	"backend/pkg/token"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthHandler struct {
	log          *zap.Logger
	repo         repo.UserRepository
	cacheManager *cache.Manager

	token *token.JWTtoken
	hash  hash.Hash
}

func NewAuthHandler(log *zap.Logger, repo repo.UserRepository, cacheManager *cache.Manager, token *token.JWTtoken, hash hash.Hash) *AuthHandler {
	return &AuthHandler{
		log:          log,
		repo:         repo,
		cacheManager: cacheManager,
		token:        token,
		hash:         hash,
	}
}

type loginRequest struct {
	Mail     string `json:"mail" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type registerRequest struct {
	Mail     string `json:"mail" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type sessionSubject struct {
	UserID string `json:"user_id"`
	Group  string `json:"group"`
}

func (a *AuthHandler) PostLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req loginRequest

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect bind: %w", err))
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect data: %w", err))
		}

		user, err := a.repo.Login(c.Request().Context(), req.Mail)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("user not found: %w", err))
		}

		verified, err := a.hash.Verify(req.Password, user.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("compare password error: %w", err))
		}

		if !verified {
			return echo.NewHTTPError(http.StatusUnauthorized, "incorrect password")
		}

		subject, err := json.Marshal(sessionSubject{
			UserID: user.ID,
			Group:  user.GroupAlias,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("marshal subject error: %w", err))
		}

		sessionID := uuid.New().String()

		if err := cache.SetWithTTL(c.Request().Context(), a.cacheManager, cache.SessionKey, sessionID, string(subject), a.token.ExpireAt); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("set session error: %w", err))
		}

		t, err := a.token.GenerateToken(sessionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("generate token error: %w", err))
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			SameSite: http.SameSiteLaxMode,
			Value:    string(t),
			Expires:  time.Now().Add(a.token.ExpireAt),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		})

		return c.NoContent(http.StatusOK)
	}
}

func (a *AuthHandler) PostRegister() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req registerRequest

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect bind: %w", err))
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect data: %w", err))
		}

		hashedPassword, err := a.hash.Hash(req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("hash password error: %w", err))
		}

		user := &domain.UserRegister{
			Mail:     req.Mail,
			Password: hashedPassword,
		}

		userData, err := a.repo.Register(c.Request().Context(), user)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return echo.NewHTTPError(http.StatusConflict, "user already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("register error: %w", err))
		}

		subject, err := json.Marshal(sessionSubject{
			UserID: userData.ID,
			Group:  userData.GroupAlias,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("marshal subject error: %w", err))
		}

		sessionID := uuid.New().String()

		if err := cache.SetWithTTL(c.Request().Context(), a.cacheManager, cache.SessionKey, sessionID, string(subject), a.token.ExpireAt); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("set session error: %w", err))
		}

		t, err := a.token.GenerateToken(sessionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("generate token error: %w", err))
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			SameSite: http.SameSiteLaxMode,
			Value:    string(t),
			Expires:  time.Now().Add(a.token.ExpireAt),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		})

		return c.NoContent(http.StatusOK)
	}
}

func (a *AuthHandler) PostLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cookie not found")
		}

		if err := cache.Delete(c.Request().Context(), a.cacheManager, cache.SessionKey, cookie.Value); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("delete session error: %w", err))
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			SameSite: http.SameSiteLaxMode,
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		})

		return c.NoContent(http.StatusOK)
	}
}
