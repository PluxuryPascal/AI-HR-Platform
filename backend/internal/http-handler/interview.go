package handler

import (
	"backend/internal/response"
	"backend/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type InterviewHandler struct {
	log *zap.Logger
	uc  usecase.InterviewUseCase
}

func NewInterviewHandler(log *zap.Logger, uc usecase.InterviewUseCase) *InterviewHandler {
	return &InterviewHandler{
		log: log,
		uc:  uc,
	}
}

func (h *InterviewHandler) GenerateQuestions() echo.HandlerFunc {
	return func(c echo.Context) error {
		candidateID := c.Param("id")
		if candidateID == "" {
			return response.Error(c, http.StatusBadRequest, "id is required")
		}

		var req struct {
			Locale string `json:"locale" validate:"required"`
		}
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid request")
		}

		teamID, ok := c.Get("team_id").(string)
		if !ok {
			return response.Error(c, http.StatusUnauthorized, "team_id not found in context")
		}

		h.log.Info("GenerateQuestions handler called", zap.String("candidate_id", candidateID))

		result, err := h.uc.GenerateQuestions(c.Request().Context(), candidateID, teamID, req.Locale)
		if err != nil {
			h.log.Error("failed to generate questions", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, err.Error())
		}

		return response.OK(c, result)
	}
}

func (h *InterviewHandler) GetQuestions() echo.HandlerFunc {
	return func(c echo.Context) error {
		candidateID := c.Param("id")
		if candidateID == "" {
			return response.Error(c, http.StatusBadRequest, "id is required")
		}

		result, err := h.uc.GetQuestions(c.Request().Context(), candidateID)
		if err != nil {
			h.log.Debug("interview guide not found for candidate", zap.String("id", candidateID))
			return response.Error(c, http.StatusNotFound, "interview guide not found")
		}

		return response.OK(c, result)
	}
}
