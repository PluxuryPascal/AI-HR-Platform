package domain

import (
	"time"
)

type EmailType string

const (
	EmailRejection       EmailType = "rejection"
	EmailInterviewInvite EmailType = "interview_invite"
)

type Communication struct {
	ID                string     `json:"id" db:"id"`
	CandidateID       string     `json:"candidate_id" db:"candidate_id"`
	GeneratedByUserID string     `json:"generated_by_user_id" db:"generated_by_user_id"`
	Type              EmailType  `json:"type" db:"type"`
	Subject           string     `json:"subject" db:"subject"`
	Body              string     `json:"body" db:"body"`
	SentAt            *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}
