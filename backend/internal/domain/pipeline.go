package domain

import (
	"time"
)

// PipelineStage represents hiring.t_pipeline_stages
type PipelineStage struct {
	ID         string  `json:"id" db:"id"`
	JobID      *string `json:"job_id,omitempty" db:"job_id"` // NULL for templates
	TeamID     string  `json:"team_id" db:"team_id"`
	Code       string  `json:"code" db:"code"`
	Title      string  `json:"title" db:"title"`
	Position   int     `json:"position" db:"position"`
	IsTerminal bool    `json:"is_terminal" db:"is_terminal"`
	Color      *string `json:"color,omitempty" db:"color"`
}

// CandidateStage represents hiring.t_candidate_stages (Current position)
type CandidateStage struct {
	CandidateID string    `json:"candidate_id" db:"candidate_id"`
	StageID     string    `json:"stage_id" db:"stage_id"`
	Position    float64   `json:"position" db:"position"`
	MovedAt     time.Time `json:"moved_at" db:"moved_at"`
}

// MoveCandidateParams for usecase/repo
type MoveCandidateParams struct {
	CandidateID string  `json:"candidate_id"`
	ToStageID   string  `json:"to_stage_id"`
	NewPosition float64 `json:"new_position"`
	ChangedBy   string  `json:"changed_by"`
}

// CreateStageParams
type CreateStageParams struct {
	JobID      *string `json:"job_id,omitempty"`
	TeamID     string  `json:"team_id"`
	Code       string  `json:"code"`
	Title      string  `json:"title"`
	Position   int     `json:"position"`
	IsTerminal bool    `json:"is_terminal"`
	Color      *string `json:"color,omitempty"`
}
