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

type PipelineHandler struct {
	log        *zap.Logger
	pipelineUC usecase.PipelineUseCase
}

func NewPipelineHandler(log *zap.Logger, uc usecase.PipelineUseCase) *PipelineHandler {
	return &PipelineHandler{
		log:        log,
		pipelineUC: uc,
	}
}

func (h *PipelineHandler) PostStage() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")
		
		var req domain.CreateStageParams
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		
		req.JobID = &jobID
		req.TeamID = c.Get("team_id").(string)
		req.ActorID = c.Get("id").(string)

		stage, err := h.pipelineUC.CreateStage(c.Request().Context(), req)
		if err != nil {
			h.log.Error("create stage error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.Created(c, stage)
	}
}

func (h *PipelineHandler) PatchStage() echo.HandlerFunc {
	return func(c echo.Context) error {
		stageID := c.Param("stage_id")
		
		var req domain.PipelineStage
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		req.ID = stageID
		req.TeamID = c.Get("team_id").(string)

		actorID := c.Get("id").(string)

		if err := h.pipelineUC.UpdateStage(c.Request().Context(), &req, actorID); err != nil {
			h.log.Error("update stage error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]bool{"success": true})
	}
}

func (h *PipelineHandler) DeleteStage() echo.HandlerFunc {
	return func(c echo.Context) error {
		stageID := c.Param("stage_id")
		actorID := c.Get("id").(string)
		teamID := c.Get("team_id").(string)

		if err := h.pipelineUC.DeleteStage(c.Request().Context(), stageID, actorID, teamID); err != nil {
			h.log.Error("delete stage error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.NoContent(c, http.StatusNoContent)
	}
}

func (h *PipelineHandler) GetStages() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")

		stages, err := h.pipelineUC.GetStagesByJobID(c.Request().Context(), jobID)
		if err != nil {
			h.log.Error("get stages error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, stages)
	}
}
