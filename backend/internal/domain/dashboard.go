package domain

import (
	"time"

	"github.com/google/uuid"
)

type DashboardStats struct {
	TotalCandidates    int     `json:"total_candidates"`
	TotalCandidatesDelta int   `json:"total_candidates_delta"`    // computed later if needed

	ActiveJobs         int     `json:"active_jobs"`
	ActiveJobsDelta    int     `json:"active_jobs_delta"`

	UpcomingInterviews int     `json:"upcoming_interviews"`
	InterviewsDelta    int     `json:"interviews_delta"`

	AvgTimeToHireDays  float64 `json:"avg_time_to_hire_days"`
	AvgTimeToHireDelta float64 `json:"avg_time_to_hire_delta"`
}

type ChartDataPoint struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type ActivityLogEntry struct {
	LogID      uuid.UUID  `json:"id"`
	ActorType  ActorType  `json:"actor_type"`
	ActorID    *uuid.UUID `json:"actor_id"`
	ActionCode string     `json:"action_code"`
	TargetID   *uuid.UUID `json:"target_id"`
	CreatedAt  time.Time  `json:"created_at"`

	// Enriched fields from simple JOINs
	CandidateFirstName *string `json:"candidate_first_name,omitempty"`
	CandidateLastName  *string `json:"candidate_last_name,omitempty"`
	JobTitle           *string `json:"job_title,omitempty"`

	// Enriched fields via gRPC
	ActorName  string `json:"actor_name"`
	MatchScore *int32 `json:"match_score,omitempty"`
}

type DashboardDynamicsRequest struct {
	StartDate string `query:"start_date"` // YYYY-MM-DD
	EndDate   string `query:"end_date"`
	Type      string `query:"type"`       // "daily" or "monthly"
}
