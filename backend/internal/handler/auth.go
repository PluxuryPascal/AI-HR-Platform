package handler

import (
	"backend/internal/repo"
	"backend/pkg/hash"
	"backend/pkg/token"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthHandler struct {
	log   *zap.Logger
	repo  repo.UserRepository
	token *token.JWTtoken
	hash  hash.Hash
}

func NewAuthHandler(log *zap.Logger, repo repo.UserRepository, token *token.JWTtoken, hash hash.Hash) *AuthHandler {
	return &AuthHandler{
		log:   log,
		repo:  repo,
		token: token,
		hash:  hash,
	}
}

type loginRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type registerRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type sessionSubject struct {
	UserID string `json:"user_id"`
	Group  string `json:"group"`
}

func (h *AuthHandler) PostLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req loginRequest

		if err := c.Bind(&req); err != nil {
			return fmt.Errorf("incorrect bind: %w", err)
		}

		if err := c.Validate(&req); err != nil {
			return fmt.Errorf("incorrect data: %w", err)
		}

		return nil
	}
}

func (h *AuthHandler) PostRegister() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req registerRequest

		if err := c.Bind(&req); err != nil {
			return fmt.Errorf("incorrect bind: %w", err)
		}

		if err := c.Validate(&req); err != nil {
			return fmt.Errorf("incorrect data: %w", err)
		}

		return nil
	}
}
