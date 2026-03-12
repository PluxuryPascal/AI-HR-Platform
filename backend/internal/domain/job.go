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
	SortID   *string `json:"sort_id"`
	SortDesc *bool   `json:"sort_desc"`
}

type DateFilter struct {
	Type     string     `json:"type"` // "before", "after", "between"
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
}

type JobFilter struct {
	Title          *string         `json:"title"`
	DepartmentName *string         `json:"department_name"`
	Status         *JobStatus      `json:"status"`
	WorkFormat     *WorkFormatType `json:"work_format"`
	DateFilter     *DateFilter     `json:"date_filter"`
	Sort           *SortParams     `json:"sort"`
}

type JobsDTO struct {
	Total int   `json:"total" db:"total"`
	Jobs  []Job `json:"jobs" db:"jobs"`
}
