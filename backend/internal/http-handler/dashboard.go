package handler

import (
	"backend/internal/domain"
	"backend/internal/response"
	"backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	log     *zap.Logger
	usecase usecase.DashboardUseCase
}

func NewDashboardHandler(log *zap.Logger, uc usecase.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{
		log:     log,
		usecase: uc,
	}
}

func (h *DashboardHandler) GetStats() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Get("team_id").(string)

		stats, err := h.usecase.GetDashboardStats(c.Request().Context(), teamID)
		if err != nil {
			h.log.Error("get dashboard stats error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, stats)
	}
}

func (h *DashboardHandler) GetApplicationDynamics() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Get("team_id").(string)

		var req domain.DashboardDynamicsRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid request format")
		}

		dynamics, err := h.usecase.GetApplicationDynamics(c.Request().Context(), req, teamID)
		if err != nil {
			h.log.Error("get dashboard application dynamics error", zap.Error(err))
			return response.Error(c, http.StatusBadRequest, err.Error())
		}

		return response.OK(c, dynamics)
	}
}

func (h *DashboardHandler) GetRecentActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Get("team_id").(string)

		limitStr := c.QueryParam("limit")
		limit := 10
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		activity, err := h.usecase.GetRecentActivity(c.Request().Context(), teamID, limit)
		if err != nil {
			h.log.Error("get dashboard recent activity error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, activity)
	}
}
