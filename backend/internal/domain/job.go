package domain

import (
	"time"
)

type JobStatus string

const (
	JobStatusDraft     JobStatus = "status_draft"
	JobStatusPublished JobStatus = "status_published"
	JobStatusClosed    JobStatus = "status_closed"
	JobStatusArchived  JobStatus = "status_archived"
)

type WorkFormatType string

const (
	WorkFormatRemote WorkFormatType = "remote"
	WorkFormatOffice WorkFormatType = "office"
	WorkFormatHybrid WorkFormatType = "hybrid"
)

// Job represents the hiring.t_jobs table
type Job struct {
	ID                    string         `json:"id" db:"id"`
	TeamID                string         `json:"team_id" db:"team_id"`
	Title                 string         `json:"title" db:"title"`
	Department            *string        `json:"department,omitempty" db:"department"`
	WorkFormat            WorkFormatType `json:"work_format" db:"work_format"`
	Description           *string        `json:"description,omitempty" db:"description"`
	ExtractedRequirements []byte         `json:"extracted_requirements,omitempty" db:"extracted_requirements"`
	Status                JobStatus      `json:"status" db:"status"`
	SalaryMin             *int           `json:"salary_min,omitempty" db:"salary_min"`
	SalaryMax             *int           `json:"salary_max,omitempty" db:"salary_max"`
	Currency              string         `json:"currency" db:"currency"`
	CreatedBy             string         `json:"created_by" db:"created_by"`
	CreatedAt             time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at" db:"updated_at"`
}
