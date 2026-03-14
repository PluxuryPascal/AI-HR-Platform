package domain

import (
	"time"
)

// EmailType represents the email_type ENUM
type EmailType string

const (
	EmailRejection       EmailType = "rejection"
	EmailInterviewInvite EmailType = "interview_invite"
)

// Communication represents the ai_engine.t_communications table
type Communication struct {
	ID                string    `json:"id" db:"id"`
	CandidateID       string    `json:"candidate_id" db:"candidate_id"`
	GeneratedByUserID string    `json:"generated_by_user_id" db:"generated_by_user_id"`
	Type              EmailType `json:"type" db:"type"`
	Content           string    `json:"content" db:"content"`
	SentAt            *time.Time `json:"sent_at,omitempty" db:"sent_at"`
}
