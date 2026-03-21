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
	DepartmentID          *string        `json:"department_id,omitempty" db:"department_id"`
	DepartmentName        *string        `json:"department_name,omitempty" db:"department_name"`
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

type SortParams struct {
	SortID   *string `json:"sort_id" validate:"omitempty"`
	SortDesc *bool   `json:"sort_desc" validate:"omitempty"`
}

type DateFilter struct {
	Type     string     `json:"type" validate:"omitempty,oneof=before after between"` // "before", "after", "between"
	DateFrom *time.Time `json:"date_from" validate:"required_if=Type after,required_if=Type between"`
	DateTo   *time.Time `json:"date_to" validate:"required_if=Type before,required_if=Type between"`
}

type JobFilter struct {
	Title          *string         `json:"title" validate:"omitempty"`
	DepartmentName *string         `json:"department_name" validate:"omitempty"`
	Status         *JobStatus      `json:"status" validate:"omitempty,oneof=status_draft status_published status_closed status_archived"`
	WorkFormat     *WorkFormatType `json:"work_format" validate:"omitempty,oneof=remote office hybrid"`
	DateFilter     *DateFilter     `json:"date_filter" validate:"omitempty"`
	Sort           *SortParams     `json:"sort" validate:"omitempty"`
	AllowedUserID *string         `json:"-"`
}

type JobsDTO struct {
	Total int   `json:"total" db:"total"`
	Jobs  []Job `json:"jobs" db:"jobs"`
}
