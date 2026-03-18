package domain

import "time"

// StageHistoryEntry represents a row from t_candidate_stage_history
// enriched with stage titles via JOIN.
type StageHistoryEntry struct {
	ID            string    `json:"id" db:"id"`
	CandidateID   string    `json:"candidate_id" db:"candidate_id"`
	FromStageID   *string   `json:"from_stage_id,omitempty" db:"from_stage_id"`
	FromStageTitle *string  `json:"from_stage_title,omitempty" db:"from_stage_title"`
	ToStageID     string    `json:"to_stage_id" db:"to_stage_id"`
	ToStageTitle  string    `json:"to_stage_title" db:"to_stage_title"`
	ChangedBy     string    `json:"changed_by" db:"changed_by"`
	ChangedAt     time.Time `json:"changed_at" db:"changed_at"`
}

// JobAccess represents hiring.t_job_access
type JobAccess struct {
	UserID string `json:"user_id" db:"user_id"`
	JobID  string `json:"job_id" db:"job_id"`
}
