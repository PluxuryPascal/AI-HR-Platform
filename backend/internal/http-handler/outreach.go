package handler

import (
	"backend/internal/domain"
	"backend/internal/response"
	"backend/internal/usecase"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OutreachHandler struct {
	outreachUC usecase.OutreachUseCase
}

func NewOutreachHandler(uc usecase.OutreachUseCase) *OutreachHandler {
	return &OutreachHandler{
		outreachUC: uc,
	}
}

type generateEmailRequest struct {
	Type   domain.EmailType `json:"type" validate:"required,oneof=rejection interview_invite"`
	Tone   string           `json:"tone" validate:"required,oneof=professional friendly brief"`
	Locale string           `json:"locale" validate:"required"`
}

type sendEmailRequest struct {
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

func (h *OutreachHandler) PostGenerateEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		candidateID := c.Param("candidate_id")
		userID := c.Get("id").(string)
		teamID := c.Get("team_id").(string)

		var req generateEmailRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		// Role is not strictly needed for generation if ai_engine fetches candidate, 
		// but usecase signature has it. I'll pass empty or just ignore in usecase.
		comm, err := h.outreachUC.GenerateEmail(c.Request().Context(), candidateID, userID, teamID, "", req.Type, req.Tone, req.Locale)
		if err != nil {
			return response.Error(c, http.StatusInternalServerError, err.Error())
		}

		return response.OK(c, comm)
	}
}

func (h *OutreachHandler) PostSendEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		communicationID := c.Param("id")

		var req sendEmailRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		if err := h.outreachUC.SendEmail(c.Request().Context(), communicationID, req.Subject, req.Body); err != nil {
			return response.Error(c, http.StatusInternalServerError, err.Error())
		}

		return response.OK(c, map[string]bool{"success": true})
	}
}
