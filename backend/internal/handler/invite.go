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

type postInviteRequest struct {
	Email  string    `json:"email"   validate:"required,email"`
	Role   string    `json:"role"    validate:"required"`
	JobIDs *[]string `json:"job_ids" validate:"omitempty,dive,uuid"`
}

type acceptInviteRequest struct {
	Token     string `json:"token"      validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Password  string `json:"password"   validate:"required,min=8"`
}

func (i *InviteHandler) GetValidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := c.QueryParam("token")
		if tokenStr == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "token is required")
		}

		payload, err := i.usecase.ValidateInvite(c.Request().Context(), tokenStr)
		if err != nil {
			if errors.Is(err, usecase.ErrInviteNotFound) || errors.Is(err, usecase.ErrInviteExpired) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("validate error: %w", err))
		}

		return c.JSON(http.StatusOK, payload)
	}
}

func (i *InviteHandler) PostInvite() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req postInviteRequest

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect bind: %w", err))
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("incorrect data: %w", err))
		}

		cookie, err := c.Cookie("access_token")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cookie not found")
		}

		if err := i.usecase.InviteUser(c.Request().Context(), cookie.Value, domain.CreateInviteParams{
			Email:  req.Email,
			Role:   req.Role,
			JobIDs: req.JobIDs,
		}); err != nil {
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

		token, expireAt, err := i.usecase.AcceptInvite(c.Request().Context(), input)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrUserAlreadyExists):
				return echo.NewHTTPError(http.StatusConflict, "user already exists")
			case errors.Is(err, usecase.ErrInviteNotFound), errors.Is(err, usecase.ErrInviteExpired):
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			default:
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("register error: %w", err))
			}
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
