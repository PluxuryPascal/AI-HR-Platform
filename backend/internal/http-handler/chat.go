package handler

import (
	"backend/internal/domain"
	"backend/internal/middleware"
	"backend/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ChatHandler struct {
	log     *zap.Logger
	useCase usecase.ChatUseCase
}

func NewChatHandler(log *zap.Logger, useCase usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		log:     log,
		useCase: useCase,
	}
}

func (h *ChatHandler) PostChat() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		teamID := middleware.GetTeamID(c)
		userID := middleware.GetUserID(c)

		var req domain.ChatRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		answer, sessionID, err := h.useCase.Chat(ctx, teamID, userID, req)
		if err != nil {
			h.log.Error("chat error", zap.Error(err), zap.String("team_id", teamID))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate chat response")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"answer":     answer,
			"session_id": sessionID,
		})
	}
}

func (h *ChatHandler) GetChatSessions() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		teamID := middleware.GetTeamID(c)
		userID := middleware.GetUserID(c)

		sessions, err := h.useCase.ListSessions(ctx, teamID, userID)
		if err != nil {
			h.log.Error("list chat sessions error", zap.Error(err), zap.String("team_id", teamID))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to list chat sessions")
		}

		return c.JSON(http.StatusOK, sessions)
	}
}

func (h *ChatHandler) GetChatHistory() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		sessionID := c.Param("id")

		history, err := h.useCase.GetHistory(ctx, sessionID)
		if err != nil {
			h.log.Error("get chat history error", zap.Error(err), zap.String("session_id", sessionID))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get chat history")
		}

		return c.JSON(http.StatusOK, history)
	}
}
