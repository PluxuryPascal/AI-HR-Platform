package handler

import (
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
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "id is required"})
		}

		var req struct {
			Locale string `json:"locale" validate:"required"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		teamID, ok := c.Get("team_id").(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "team_id not found in context"})
		}

		h.log.Info("GenerateQuestions handler called", zap.String("candidate_id", candidateID))

		result, err := h.uc.GenerateQuestions(c.Request().Context(), candidateID, teamID, req.Locale)
		if err != nil {
			h.log.Error("failed to generate questions", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, result)
	}
}
