package handler

import (
	"backend/internal/domain"
	"backend/internal/response"
	"backend/internal/usecase"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type DepartmentHandler struct {
	log    *zap.Logger
	deptUC usecase.DepartmentUseCase
}

func NewDepartmentHandler(log *zap.Logger, deptUC usecase.DepartmentUseCase) *DepartmentHandler {
	return &DepartmentHandler{
		log:    log,
		deptUC: deptUC,
	}
}

func (h *DepartmentHandler) PostDepartment() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req domain.CreateDepartmentParams
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		
		req.TeamID = c.Get("team_id").(string)

		dept, err := h.deptUC.CreateDepartment(c.Request().Context(), &req)
		if err != nil {
			h.log.Error("create department error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.Created(c, dept)
	}
}

func (h *DepartmentHandler) GetDepartments() echo.HandlerFunc {
	return func(c echo.Context) error {
		teamID := c.Get("team_id").(string)

		depts, err := h.deptUC.GetDepartmentsByTeam(c.Request().Context(), teamID)
		if err != nil {
			h.log.Error("get departments error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, depts)
	}
}

func (h *DepartmentHandler) DeleteDepartment() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		if err := h.deptUC.DeleteDepartment(c.Request().Context(), id); err != nil {
			h.log.Error("delete department error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.NoContent(c, http.StatusNoContent)
	}
}
