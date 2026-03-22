package handler

import (
	"backend/internal/middleware"
	"backend/internal/response"
	"backend/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AiHandler struct {
	usecase usecase.AiUseCase
}

func NewAiHandler(uc usecase.AiUseCase) *AiHandler {
	return &AiHandler{
		usecase: uc,
	}
}

type parseJobRequest struct {
	RawText string `json:"raw_text" validate:"required"`
	Locale  string `json:"locale"   validate:"required"`
}

func (h *AiHandler) ParseJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req parseJobRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid request")
		}

		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, "validation failed")
		}

		teamID := middleware.GetTeamID(c)

		result, err := h.usecase.ParseJobDescription(c.Request().Context(), req.RawText, req.Locale, teamID)
		if err != nil {
			return response.Error(c, http.StatusInternalServerError, err.Error())
		}

		return response.OK(c, result)
	}
}
