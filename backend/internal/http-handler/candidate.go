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

type CandidateHandler struct {
	log     *zap.Logger
	usecase usecase.CandidateUseCase
}

func NewCandidateHandler(log *zap.Logger, uc usecase.CandidateUseCase) *CandidateHandler {
	return &CandidateHandler{
		log:     log,
		usecase: uc,
	}
}

type createCandidateRequest struct {
	InitialStageID string  `json:"initial_stage_id" validate:"required,uuid4"`
	FirstName      *string `json:"first_name" validate:"omitempty,min=1,max=64"`
	LastName       *string `json:"last_name" validate:"omitempty,min=1,max=64"`
	Email          *string `json:"email" validate:"omitempty,email"`
	Location       *string `json:"location" validate:"omitempty,max=128"`
}

type listCandidatesRequest struct {
	Pagination domain.Pagination      `json:"pagination" validate:"required"`
	Filter     domain.CandidateFilter `json:"filter" validate:"omitempty"`
}

type candidateDetailResponse struct {
	Candidate *domain.Candidate        `json:"candidate"`
	Profile   *domain.CandidateProfile `json:"profile"`
	StageID   *string                  `json:"stage_id"`
}

type updateCandidateRequest struct {
	FirstName *string  `json:"first_name" validate:"omitempty,min=1,max=64"`
	LastName  *string  `json:"last_name" validate:"omitempty,min=1,max=64"`
	Email     *string  `json:"email" validate:"omitempty,email"`
	Location  *string  `json:"location" validate:"omitempty,max=128"`
	Skills    []string `json:"skills" validate:"omitempty,dive,min=1"`
}

type moveCandidateRequest struct {
	ToStageID   string  `json:"to_stage_id" validate:"required,uuid4"`
	NewPosition float64 `json:"new_position" validate:"required,min=0"`
}

func (h *CandidateHandler) PostCandidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")
		var req createCandidateRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		candidate := &domain.Candidate{
			JobID:     jobID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Location:  req.Location,
		}
		profile := &domain.CandidateProfile{}

		if err := h.usecase.CreateCandidate(c.Request().Context(), candidate, profile, req.InitialStageID); err != nil {
			h.log.Error("create candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.Created(c, candidate)
	}
}

func (h *CandidateHandler) PostCandidateList() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")
		var req listCandidatesRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		offset := (req.Pagination.Page - 1) * req.Pagination.PerPage
		dto, err := h.usecase.GetCandidatesByJob(c.Request().Context(), jobID, offset, req.Pagination.PerPage, req.Filter)
		if err != nil {
			h.log.Error("get candidates error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.Paginated(c, dto.Candidates, dto.Total, req.Pagination.Page, req.Pagination.PerPage)
	}
}

func (h *CandidateHandler) GetCandidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		candidate, profile, stageID, err := h.usecase.GetCandidateByID(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, usecase.ErrCandidateNotFound) {
				return response.Error(c, http.StatusNotFound, "candidate not found")
			}
			h.log.Error("get candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, candidateDetailResponse{
			Candidate: candidate,
			Profile:   profile,
			StageID:   stageID,
		})
	}
}

func (h *CandidateHandler) PatchCandidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req updateCandidateRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		candidate, _, _, err := h.usecase.GetCandidateByID(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, usecase.ErrCandidateNotFound) {
				return response.Error(c, http.StatusNotFound, "candidate not found")
			}
			h.log.Error("get candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		if req.FirstName != nil {
			candidate.FirstName = req.FirstName
		}
		if req.LastName != nil {
			candidate.LastName = req.LastName
		}
		if req.Email != nil {
			candidate.Email = req.Email
		}
		if req.Location != nil {
			candidate.Location = req.Location
		}
		if req.Skills != nil {
			candidate.Skills = req.Skills
		}

		if err := h.usecase.UpdateCandidate(c.Request().Context(), candidate); err != nil {
			h.log.Error("update candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, candidate)
	}
}

func (h *CandidateHandler) PostCandidateMove() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		userID := c.Get("id").(string)
		var req moveCandidateRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		params := domain.MoveCandidateParams{
			CandidateID: id,
			ToStageID:   req.ToStageID,
			NewPosition: req.NewPosition,
			ChangedBy:   userID,
		}

		if err := h.usecase.MoveCandidate(c.Request().Context(), params); err != nil {
			h.log.Error("move candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]bool{"success": true})
	}
}

func (h *CandidateHandler) DeleteCandidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := h.usecase.DeleteCandidate(c.Request().Context(), id); err != nil {
			h.log.Error("delete candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.NoContent(c, http.StatusNoContent)
	}
}
