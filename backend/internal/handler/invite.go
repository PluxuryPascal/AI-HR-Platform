package handler

import (
	"backend/internal/domain"
	"backend/internal/usecase"
	"backend/pkg/config"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type InviteHandler struct {
	cfg     *config.Server
	log     *zap.Logger
	usecase usecase.InviteUseCase
}

func NewInviteHandler(cfg *config.Server, log *zap.Logger, usecase usecase.InviteUseCase) *InviteHandler {
	return &InviteHandler{
		cfg:     cfg,
		log:     log,
		usecase: usecase,
	}
}

type acceptInviteRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (i *InviteHandler) PostValidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req domain.InvitePayload

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect bind: %w", err))
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect data: %w", err))
		}

		if err := i.usecase.InviteUser(c.Request().Context(), req); err != nil {
			if errors.Is(err, usecase.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, "user already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invite error: %w", err))
		}

		return c.NoContent(http.StatusCreated)
	}
}

func (i *InviteHandler) PostInvite() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req domain.InvitePayload

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect bind: %w", err))
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect data: %w", err))
		}

		if err := i.usecase.InviteUser(c.Request().Context(), req); err != nil {
			if errors.Is(err, usecase.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, "user already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invite error: %w", err))
		}

		return c.NoContent(http.StatusCreated)
	}
}

func (i *InviteHandler) PostCreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req acceptInviteRequest

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect bind: %w", err))
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect data: %w", err))
		}

		input := domain.CreateUserParams{
			Token:     req.Token,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		token, expireAt, err := i.usecase.CreateUser(c.Request().Context(), input)
		if err != nil {
			if errors.Is(err, usecase.ErrUserAlreadyExists) {
				return echo.NewHTTPError(http.StatusConflict, "user already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("register error: %w", err))
		}

		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			SameSite: http.SameSiteStrictMode,
			Value:    *token,
			Expires:  time.Now().Add(expireAt),
			Path:     "/",
			Secure:   i.cfg.SecureCookie,
			HttpOnly: true,
		})

		return c.NoContent(http.StatusCreated)
	}
}
