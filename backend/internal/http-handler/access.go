package handler

import (
	"backend/internal/response"
	"backend/internal/usecase"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AccessHandler struct {
	log      *zap.Logger
	accessUC usecase.AccessUseCase
}

func NewAccessHandler(log *zap.Logger, uc usecase.AccessUseCase) *AccessHandler {
	return &AccessHandler{
		log:      log,
		accessUC: uc,
	}
}

type grantAccessRequest struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
}

func (h *AccessHandler) PostGrantAccess() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")
		var req grantAccessRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		actorID := c.Get("id").(string)
		teamID := c.Get("team_id").(string)
		if err := h.accessUC.GrantAccess(c.Request().Context(), req.UserID, jobID, actorID, teamID); err != nil {
			h.log.Error("grant access error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]bool{"success": true})
	}
}

func (h *AccessHandler) DeleteRevokeAccess() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")
		userID := c.Param("user_id")

		actorID := c.Get("id").(string)
		teamID := c.Get("team_id").(string)

		if err := h.accessUC.RevokeAccess(c.Request().Context(), userID, jobID, actorID, teamID); err != nil {
			h.log.Error("revoke access error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.NoContent(c, http.StatusNoContent)
	}
}

func (h *AccessHandler) GetJobAccessList() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")

		list, err := h.accessUC.GetAccessByJobID(c.Request().Context(), jobID)
		if err != nil {
			h.log.Error("get job access list error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, list)
	}
}
