package domain

import (
	"time"
)

// TaskStatus represents the task_status ENUM
type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskExtracting TaskStatus = "extracting"
	TaskAnalyzing  TaskStatus = "analyzing"
	TaskCompleted  TaskStatus = "completed"
	TaskFailed     TaskStatus = "failed"
)

// ProcessingTask represents the ai_engine.t_processing_tasks table
type ProcessingTask struct {
	ID              string     `json:"id" db:"id"`
	WorkflowID      *string    `json:"workflow_id,omitempty" db:"workflow_id"`
	EntityID        *string    `json:"entity_id,omitempty" db:"entity_id"`
	Status          TaskStatus `json:"status" db:"status"`
	ProgressPercent int        `json:"progress_percent" db:"progress_percent"`
	ErrorMessage    *string    `json:"error_message,omitempty" db:"error_message"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
