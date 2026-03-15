package temporal

import (
	"backend/internal/domain"
	hiringv1 "backend/internal/proto/hiring/v1"
)

type ResumePipelineInput struct {
	CandidateID   string `json:"candidate_id"`
	JobID         string `json:"job_id"`
	ResumeFileKey string `json:"resume_file_key"`
	TeamID        string `json:"team_id"`
}

type LLMParseInput struct {
	CandidateID string `json:"candidate_id"`
	JobID       string `json:"job_id"`
	ParsedText  string `json:"parsed_text"`
}

type LLMScoreInput struct {
	CandidateID     string `json:"candidate_id"`
	ParsedText      string `json:"parsed_text"`
	JobRequirements string `json:"job_requirements"`
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
