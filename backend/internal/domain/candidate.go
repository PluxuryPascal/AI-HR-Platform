package domain

import (
	"io"
	"time"
)

// Candidate represents the hiring.t_candidates table
type Candidate struct {
	ID            string                 `json:"id" db:"id"`
	JobID         string                 `json:"job_id" db:"job_id"`
	FirstName     *string                `json:"first_name,omitempty" db:"first_name"`
	LastName      *string                `json:"last_name,omitempty" db:"last_name"`
	Email         *string                `json:"email,omitempty" db:"email"`
	ResumeFileKey *string                `json:"resume_file_key,omitempty" db:"resume_file_key"`
	ParsedText    *string                `json:"parsed_text,omitempty" db:"parsed_text"`
	Location      *string                `json:"location,omitempty" db:"location"`
	Skills        []string               `json:"skills,omitempty" db:"skills"`
	ParsingStatus CandidateParsingStatus `json:"parsing_status" db:"parsing_status"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at,omitempty" db:"updated_at"`
}

type CandidateParsingStatus string

const (
	ParsingStatusPending     CandidateParsingStatus = "pending"
	ParsingStatusProcessing  CandidateParsingStatus = "processing"
	ParsingStatusNeedsReview CandidateParsingStatus = "needs_review"
	ParsingStatusCompleted   CandidateParsingStatus = "completed"
	ParsingStatusFailed      CandidateParsingStatus = "failed"
)

// CandidateProfile represents the hiring.t_candidate_profiles table (1:1 with Candidate)
type CandidateProfile struct {
	CandidateID    string     `json:"candidate_id" db:"candidate_id"`
	StructuredData []byte     `json:"structured_data,omitempty" db:"structured_data"`
	AIParsedAt     *time.Time `json:"ai_parsed_at,omitempty" db:"ai_parsed_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	MissingFields  []string   `json:"missing_fields,omitempty" db:"missing_fields"`
}

// FactorType represents the factor_type ENUM
type FactorType string

const (
	FactorPositive FactorType = "positive"
	FactorNegative FactorType = "negative"
	FactorNeutral  FactorType = "neutral"
)

// CandidateScore represents the ai_engine.t_candidate_scores table
type CandidateScore struct {
	CandidateID string     `json:"candidate_id" db:"candidate_id"`
	MatchScore  int        `json:"match_score" db:"match_score"`
	AnalyzedAt  *time.Time `json:"analyzed_at,omitempty" db:"analyzed_at"`
}

// ScoreFactor represents the ai_engine.t_score_factors table
type ScoreFactor struct {
	ID          string     `json:"id" db:"id"`
	CandidateID string     `json:"candidate_id" db:"candidate_id"`
	Type        FactorType `json:"type" db:"type"`
	Description string     `json:"description" db:"description"`
	Impact      int        `json:"impact" db:"impact"`
}

// ResumeEmbedding represents the ai_engine.t_resume_embeddings table
type ResumeEmbedding struct {
	ID          string    `json:"id" db:"id"`
	TeamID      string    `json:"team_id" db:"team_id"`
	CandidateID string    `json:"candidate_id" db:"candidate_id"`
	ChunkText   string    `json:"chunk_text" db:"chunk_text"`
	Embedding   []float32 `json:"embedding" db:"embedding"`
}

type CandidateFilter struct {
	FirstName      *string     `json:"first_name"`
	LastName       *string     `json:"last_name"`
	Email          *string     `json:"email"`
	CurrentStageID *string     `json:"current_stage_id"`
	DateFilter     *DateFilter `json:"date_filter"`
	Sort           *SortParams `json:"sort"`
}

type CandidatesDTO struct {
	Total      int         `json:"total" db:"total"`
	Candidates []Candidate `json:"candidates" db:"candidates"`
}

type CreateCandidateParams struct {
	JobID    string
	Filename string
	File     io.Reader
}

type AIParsingResult struct {
	CandidateID    string
	JobID          string
	ParseStatus    CandidateParsingStatus
	FirstName      *string
	LastName       *string
	Email          *string
	Location       *string
	Skills         []string
	ParsedText     *string
	StructuredData []byte
	InitialStageID *string
	MissingFields  []string
}
