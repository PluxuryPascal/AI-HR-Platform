package domain

import (
	"time"
)

// Candidate represents the hiring.t_candidates table
type Candidate struct {
	ID            string     `json:"id" db:"id"`
	JobID         string     `json:"job_id" db:"job_id"`
	FirstName     *string    `json:"first_name,omitempty" db:"first_name"`
	LastName      *string    `json:"last_name,omitempty" db:"last_name"`
	Email         *string    `json:"email,omitempty" db:"email"`
	ResumeFileKey *string    `json:"resume_file_key,omitempty" db:"resume_file_key"`
	ParsedText    *string    `json:"parsed_text,omitempty" db:"parsed_text"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// CandidateProfile represents the hiring.t_candidate_profiles table (1:1 with Candidate)
type CandidateProfile struct {
	CandidateID    string     `json:"candidate_id" db:"candidate_id"`
	StructuredData []byte     `json:"structured_data,omitempty" db:"structured_data"`
	AIParsedAt     *time.Time `json:"ai_parsed_at,omitempty" db:"ai_parsed_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
