package handler

import (
	"backend/internal/usecase"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthHandler struct {
	log     *zap.Logger
	usecase usecase.AuthUseCase
}

func NewAuthHandler(log *zap.Logger, usecase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		log:     log,
		usecase: usecase,
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

		token, expireAt, err := a.usecase.Login(c.Request().Context(), req.Mail, req.Password)
		if err != nil {
			if errors.Is(err, usecase.ErrInvalidCredentials) {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("login error: %w", err))
		}

		if token == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "token generation failed")
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			SameSite: http.SameSiteLaxMode,
			Value:    *token,
			Expires:  time.Now().Add(expireAt),
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

		token, expireAt, err := a.usecase.Register(c.Request().Context(), req.Mail, req.Password)
		if err != nil {
			if errors.Is(err, usecase.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, "user already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("register error: %w", err))
		}

		if token == nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "token generation failed")
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			SameSite: http.SameSiteLaxMode,
			Value:    *token,
			Expires:  time.Now().Add(expireAt),
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

		if err := a.usecase.Logout(c.Request().Context(), cookie.Value); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("logout error: %w", err))
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
