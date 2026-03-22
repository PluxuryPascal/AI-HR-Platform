package activity

import (
	"backend/internal/domain"
	hiringv1 "backend/internal/proto/hiring/v1"
)

type ResumePipelineInput struct {
	CandidateID   string `json:"candidate_id"`
	JobID         string `json:"job_id"`
	ResumeFileKey string `json:"resume_file_key"`
	TeamID        string `json:"team_id"`
	Locale        string `json:"locale"`
}

type LLMParseInput struct {
	TeamID      string `json:"team_id"`
	CandidateID string `json:"candidate_id"`
	JobID       string `json:"job_id"`
	ParsedText  string `json:"parsed_text"`
	Locale      string `json:"locale"`
}

type LLMScoreInput struct {
	CandidateID     string `json:"candidate_id"`
	TeamID          string `json:"team_id"`
	ParsedText      string `json:"parsed_text"`
	JobRequirements string `json:"job_requirements"`
	Locale          string `json:"locale"`
}

type LLMEmbedInput struct {
	CandidateID string `json:"candidate_id"`
	TeamID      string `json:"team_id"`
	ParsedText  string `json:"parsed_text"`
}

type GRPCCallbackInput struct {
	CandidateID string                 `json:"candidate_id"`
	JobID       string                 `json:"job_id"`
	Status      hiringv1.ParsingStatus `json:"status"`
	ParsedText  string                 `json:"parsed_text"`
	ParseResult domain.ParseResult     `json:"parse_result"`
}

type EmailGenerateInput struct {
	CandidateID        string   `json:"candidate_id"`
	TeamID             string   `json:"team_id"`
	GeneratedByUserID  string   `json:"generated_by_user_id"`
	CandidateFirstName string   `json:"candidate_first_name"`
	CandidateLastName  string   `json:"candidate_last_name"`
	CandidateEmail     string   `json:"candidate_email"`
	Role               string   `json:"role"`
	Skills             []string `json:"skills"`
	MatchScore         int      `json:"match_score"`
	EmailType          string   `json:"email_type"`
	Tone               string   `json:"tone"`
	Locale             string   `json:"locale"`
	RecruiterName      string   `json:"recruiter_name"`
	CompanyName        string   `json:"company_name"`
}

type EmailGenerateOutput struct {
	CommunicationID string `json:"communication_id"`
	Subject         string `json:"subject"`
	Body            string `json:"body"`
}

type CandidateCompareCandidate struct {
	ID         string   `json:"id"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Role       string   `json:"role"`
	MatchScore int      `json:"match_score"`
	Skills     []string `json:"skills"`
	ParsedText string   `json:"parsed_text"`
}

type CandidateCompareInput struct {
	TeamID          string                      `json:"team_id"`
	Candidates      []CandidateCompareCandidate `json:"candidates"`
	JobRequirements string                      `json:"job_requirements"`
	Locale          string                      `json:"locale"`
}

type CompareResultEntry struct {
	Experience  string   `json:"experience"`
	Skills      []string `json:"skills"`
	SalaryRange string   `json:"salary_range"`
	Risks       string   `json:"risks"`
	Summary     string   `json:"summary"`
}

type CandidateCompareOutput map[string]CompareResultEntry

type JobParseInput struct {
	TeamID  string `json:"team_id"`
	RawText string `json:"raw_text"`
	Locale  string `json:"locale"`
}

type JobParseOutput struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Requirements []string `json:"requirements"`
	WorkFormat   string   `json:"work_format"`
	SalaryMin    int      `json:"salary_min"`
	SalaryMax    int      `json:"salary_max"`
	Currency     string   `json:"currency"`
}

type ChatInput struct {
	SessionID   string               `json:"session_id"`
	TeamID      string               `json:"team_id"`
	CandidateID *string              `json:"candidate_id,omitempty"`
	Question    string               `json:"question"`
	Locale      string               `json:"locale"`
	History     []domain.ChatMessage `json:"history,omitempty"`
}

type ChatOutput struct {
	Answer string `json:"answer"`
}

type InterviewPair struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Category string `json:"category"`
}

type InterviewInput struct {
	CandidateID string `json:"candidate_id"`
	TeamID      string `json:"team_id"`
	Locale      string `json:"locale"`
}

type InterviewOutput struct {
	Questions []InterviewPair `json:"questions"`
}
