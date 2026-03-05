package domain

import (
	"time"
)

// PipelineStage represents the hiring.t_pipeline_stages table
type PipelineStage struct {
	ID         string  `json:"id" db:"id"`
	JobID      *string `json:"job_id,omitempty" db:"job_id"` // nil = team template
	TeamID     string  `json:"team_id" db:"team_id"`
	Code       string  `json:"code" db:"code"`
	Title      string  `json:"title" db:"title"`
	Position   int     `json:"position" db:"position"`
	IsTerminal bool    `json:"is_terminal" db:"is_terminal"`
	Color      *string `json:"color,omitempty" db:"color"`
}

// CandidateStage represents the hiring.t_candidate_stages table (candidate's current position)
type CandidateStage struct {
	CandidateID string    `json:"candidate_id" db:"candidate_id"`
	StageID     string    `json:"stage_id" db:"stage_id"`
	Position    float64   `json:"position" db:"position"`
	MovedAt     time.Time `json:"moved_at" db:"moved_at"`
}

// CandidateStageHistory represents the hiring.t_candidate_stage_history table
type CandidateStageHistory struct {
	ID          string    `json:"id" db:"id"`
	CandidateID string    `json:"candidate_id" db:"candidate_id"`
	FromStageID *string   `json:"from_stage_id,omitempty" db:"from_stage_id"`
	ToStageID   string    `json:"to_stage_id" db:"to_stage_id"`
	ChangedBy   string    `json:"changed_by" db:"changed_by"`
	ChangedAt   time.Time `json:"changed_at" db:"changed_at"`
}
