package handler

import (
	"backend/internal/domain"
	"backend/internal/response"
	"backend/internal/usecase"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type JobHandler struct {
	log     *zap.Logger
	usecase usecase.JobUseCase
}

func NewJobHandler(log *zap.Logger, uc usecase.JobUseCase) *JobHandler {
	return &JobHandler{
		log:     log,
		usecase: uc,
	}
}

type createJobRequest struct {
	Title        string                `json:"title" validate:"required,min=2,max=128"`
	DepartmentID *string               `json:"department_id,omitempty" validate:"omitempty,uuid4"`
	WorkFormat   domain.WorkFormatType `json:"work_format" validate:"required,oneof=remote office hybrid"`
	Description  *string               `json:"description,omitempty"`
	SalaryMin    *int                  `json:"salary_min,omitempty" validate:"omitempty,min=0"`
	SalaryMax    *int                  `json:"salary_max,omitempty" validate:"omitempty,min=0"`
	Currency     string                `json:"currency,omitempty" validate:"omitempty,len=3"`
}

type jobListRequest struct {
	Pagination domain.Pagination `json:"pagination" validate:"required"`
	Filter     domain.JobFilter  `json:"filter" validate:"omitempty"`
}

type updateJobRequest struct {
	Title        *string                `json:"title,omitempty"`
	DepartmentID *string                `json:"department_id,omitempty" validate:"omitempty,uuid4"`
	WorkFormat   *domain.WorkFormatType `json:"work_format,omitempty" validate:"omitempty,oneof=remote office hybrid"`
	Description  *string                `json:"description,omitempty"`
	SalaryMin    *int                   `json:"salary_min,omitempty" validate:"omitempty,min=0"`
	SalaryMax    *int                   `json:"salary_max,omitempty" validate:"omitempty,min=0"`
	Currency     *string                `json:"currency,omitempty" validate:"omitempty,len=3"`
}

func (h *JobHandler) PostJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req createJobRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}

		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		currency := req.Currency
		if currency == "" {
			currency = "RUB"
		}

		job := &domain.Job{
			Title:        req.Title,
			DepartmentID: req.DepartmentID,
			WorkFormat:   req.WorkFormat,
			Description:  req.Description,
			SalaryMin:    req.SalaryMin,
			SalaryMax:    req.SalaryMax,
			Currency:     currency,
			Status:       domain.JobStatusDraft,
			CreatedBy:    c.Get("id").(string),
			TeamID:       c.Get("team_id").(string),
		}

		if err := h.usecase.CreateJob(c.Request().Context(), job); err != nil {
			h.log.Error("create job error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.Created(c, job)
	}
}

func (h *JobHandler) PostJobList() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req jobListRequest

		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}

		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		teamID := c.Get("team_id").(string)
		offset := (req.Pagination.Page - 1) * req.Pagination.PerPage

		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		dto, err := h.usecase.GetJobsByTeam(c.Request().Context(), teamID, userID, role, offset, req.Pagination.PerPage, req.Filter)
		if err != nil {
			h.log.Error("get jobs error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.Paginated(c, dto.Jobs, dto.Total, req.Pagination.Page, req.Pagination.PerPage)
	}
}

func (h *JobHandler) GetJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		job, err := h.usecase.GetJobByID(c.Request().Context(), id, userID, role)
		if err != nil {
			if errors.Is(err, usecase.ErrJobNotFound) {
				return response.Error(c, http.StatusNotFound, "job not found")
			}
			h.log.Error("get job error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, job)
	}
}

func (h *JobHandler) PatchJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		var req updateJobRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}

		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		job, err := h.usecase.GetJobByID(c.Request().Context(), id, userID, role)
		if err != nil {
			if errors.Is(err, usecase.ErrJobNotFound) {
				return response.Error(c, http.StatusNotFound, "job not found")
			}
			h.log.Error("get job error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		if req.Title != nil {
			job.Title = *req.Title
		}
		if req.DepartmentID != nil {
			job.DepartmentID = req.DepartmentID
		}
		if req.WorkFormat != nil {
			job.WorkFormat = *req.WorkFormat
		}
		if req.Description != nil {
			job.Description = req.Description
		}
		if req.SalaryMin != nil {
			job.SalaryMin = req.SalaryMin
		}
		if req.SalaryMax != nil {
			job.SalaryMax = req.SalaryMax
		}
		if req.Currency != nil {
			job.Currency = *req.Currency
		}

		actorID := c.Get("id").(string)
		if err := h.usecase.UpdateJob(c.Request().Context(), job, actorID, role); err != nil {
			h.log.Error("update job error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, job)
	}
}

func (h *JobHandler) PostJobPublish() echo.HandlerFunc {
	return h.changeStatus(domain.JobStatusPublished)
}

func (h *JobHandler) PostJobClose() echo.HandlerFunc {
	return h.changeStatus(domain.JobStatusClosed)
}

func (h *JobHandler) PostJobArchive() echo.HandlerFunc {
	return h.changeStatus(domain.JobStatusArchived)
}

func (h *JobHandler) changeStatus(newStatus domain.JobStatus) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		userID := c.Get("id").(string)
		role := c.Get("role").(string)

		job, err := h.usecase.GetJobByID(c.Request().Context(), id, userID, role)
		if err != nil {
			if errors.Is(err, usecase.ErrJobNotFound) {
				return response.Error(c, http.StatusNotFound, "job not found")
			}
			h.log.Error("get job error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		// Validate allowed transitions
		if !isValidStatusTransition(job.Status, newStatus) {
			return response.Error(c, http.StatusConflict,
				fmt.Sprintf("invalid status transition: %s → %s", job.Status, newStatus))
		}

		job.Status = newStatus
		actorID := c.Get("id").(string)
		if err := h.usecase.UpdateJob(c.Request().Context(), job, actorID, role); err != nil {
			h.log.Error("update job status error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, job)
	}
}

// isValidStatusTransition checks allowed job status transitions:
// Draft → Published, Published → Closed, Closed → Archived
func isValidStatusTransition(from, to domain.JobStatus) bool {
	transitions := map[domain.JobStatus]domain.JobStatus{
		domain.JobStatusDraft:     domain.JobStatusPublished,
		domain.JobStatusPublished: domain.JobStatusClosed,
		domain.JobStatusClosed:    domain.JobStatusArchived,
	}
	allowed, ok := transitions[from]
	return ok && allowed == to
}

func (h *JobHandler) DeleteJob() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		actorID := c.Get("id").(string)
		teamID := c.Get("team_id").(string)
		role := c.Get("role").(string)
		if err := h.usecase.DeleteJob(c.Request().Context(), id, actorID, teamID, role); err != nil {
			h.log.Error("delete job error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.NoContent(c, http.StatusNoContent)
	}
}
