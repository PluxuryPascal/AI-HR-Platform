package handler

import (
	"backend/internal/domain"
	"backend/internal/response"
	"backend/internal/usecase"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	maxResumeSize = 10 << 20 // 10 MB на один файл
)

type CandidateHandler struct {
	log         *zap.Logger
	candidateUC usecase.CandidateUseCase
}

func NewCandidateHandler(log *zap.Logger, uc usecase.CandidateUseCase) *CandidateHandler {
	return &CandidateHandler{
		log:         log,
		candidateUC: uc,
	}
}

type listCandidatesRequest struct {
	Pagination domain.Pagination      `json:"pagination" validate:"required"`
	Filter     domain.CandidateFilter `json:"filter" validate:"omitempty"`
}

type candidateDetailResponse struct {
	Candidate *domain.Candidate        `json:"candidate"`
	Profile   *domain.CandidateProfile `json:"profile"`
	StageID   *string                  `json:"stage_id"`
	Score     *domain.CandidateScore   `json:"score"`
	Factors   []domain.ScoreFactor     `json:"factors"`
}

type updateCandidateRequest struct {
	FirstName *string  `json:"first_name" validate:"omitempty,min=1,max=64"`
	LastName  *string  `json:"last_name" validate:"omitempty,min=1,max=64"`
	Email     *string  `json:"email" validate:"omitempty,email"`
	Location  *string  `json:"location" validate:"omitempty,max=128"`
	Skills    []string `json:"skills" validate:"omitempty,dive,min=1"`
}

type moveCandidateRequest struct {
	ToStageID   string   `json:"to_stage_id" validate:"required,uuid4"`
	NewPosition *float64 `json:"new_position" validate:"required,min=0"`
}

type uploadResumeResponse struct {
	Candidate     *domain.Candidate
	ResumeFileKey string
}

func (h *CandidateHandler) PostUploadResume() echo.HandlerFunc {
	return func(c echo.Context) error {
		jobID := c.Param("job_id")

		fh, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "field 'file' is required")
		}

		if fh.Size > maxResumeSize {
			return echo.NewHTTPError(http.StatusRequestEntityTooLarge,
				fmt.Sprintf("file too large: max %d MB", maxResumeSize>>20))
		}

		if strings.ToLower(filepath.Ext(fh.Filename)) != ".pdf" {
			return echo.NewHTTPError(http.StatusUnsupportedMediaType, "only PDF files are accepted")
		}

		file, err := fh.Open()
		if err != nil {
			h.log.Error("open uploaded file", zap.String("filename", fh.Filename), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read uploaded file")
		}
		defer file.Close()

		locale := c.FormValue("locale")
		if locale == "" {
			locale = "ru" // Default
		}

		candidate, err := h.candidateUC.CreateCandidate(c.Request().Context(), domain.CreateCandidateParams{
			JobID:    jobID,
			Filename: fh.Filename,
			File:     file,
			ActorID:  c.Get("id").(string),
			Locale:   locale,
		})
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrNoPipelineStages):
				return echo.NewHTTPError(http.StatusConflict, "job has no pipeline stages configured")
			case errors.Is(err, usecase.ErrStorageUnavailable):
				return echo.NewHTTPError(http.StatusBadGateway, "failed to upload file to storage")
			default:
				h.log.Error("upload resume", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to upload resume")
			}
		}

		return response.Created(c, uploadResumeResponse{
			Candidate:     candidate,
			ResumeFileKey: *candidate.ResumeFileKey,
		})
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
		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		dto, err := h.candidateUC.GetCandidatesByJob(c.Request().Context(), jobID, userID, role, offset, req.Pagination.PerPage, req.Filter)
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
		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		candidate, profile, stageID, score, factors, err := h.candidateUC.GetCandidateByID(c.Request().Context(), id, userID, role)
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
			Score:     score,
			Factors:   factors,
		})
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
			NewPosition: *req.NewPosition,
			ChangedBy:   userID,
		}

		role := c.Get("role").(string)
		if err := h.candidateUC.MoveCandidate(c.Request().Context(), params, role); err != nil {
			if errors.Is(err, usecase.ErrInvalidStageTransition) {
				return response.Error(c, http.StatusConflict, err.Error())
			}
			h.log.Error("move candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]bool{"success": true})
	}
}

func (h *CandidateHandler) DeleteCandidate() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		actorID := c.Get("id").(string)
		teamID := c.Get("team_id").(string)
		role := c.Get("role").(string)
		if err := h.candidateUC.DeleteCandidate(c.Request().Context(), id, actorID, teamID, role); err != nil {
			h.log.Error("delete candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.NoContent(c, http.StatusNoContent)
	}
}

func (h *CandidateHandler) GetCandidateHistory() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		entries, err := h.candidateUC.GetStageHistory(c.Request().Context(), id, userID, role)
		if err != nil {
			h.log.Error("get candidate history error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, entries)
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

		candidate := &domain.Candidate{
			ID:        id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Location:  req.Location,
			Skills:    req.Skills,
		}

		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		if err := h.candidateUC.UpdateByRecruiter(c.Request().Context(), candidate, userID, role); err != nil {
			h.log.Error("update candidate error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, candidate)
	}
}

func (h *CandidateHandler) PostBulkCandidateMove() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			CandidateIDs []string `json:"candidate_ids" validate:"required,dive,uuid4"`
			ToStageID    string   `json:"to_stage_id" validate:"required,uuid4"`
		}
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}
		if err := c.Validate(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		}

		params := domain.BulkMoveCandidateParams{
			CandidateIDs: req.CandidateIDs,
			ToStageID:    req.ToStageID,
			ChangedBy:    c.Get("id").(string),
		}

		role := c.Get("role").(string)
		if err := h.candidateUC.BulkMoveCandidates(c.Request().Context(), params, role); err != nil {
			h.log.Error("bulk move candidates error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]bool{"success": true})
	}
}

func (h *CandidateHandler) PostConfirmReview() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var req updateCandidateRequest
		if err := c.Bind(&req); err != nil {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("bind: %v", err))
		}

		candidate := &domain.Candidate{
			ID:        id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Location:  req.Location,
			Skills:    req.Skills,
		}

		userID := c.Get("id").(string)
		role := c.Get("role").(string)
		if err := h.candidateUC.ConfirmManualReview(c.Request().Context(), candidate, userID, role); err != nil {
			if errors.Is(err, usecase.ErrCandidateNotNeedsReview) {
				return response.Error(c, http.StatusConflict, err.Error())
			}
			h.log.Error("confirm review error", zap.Error(err))
			return response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return response.OK(c, map[string]bool{"success": true})
	}
}
