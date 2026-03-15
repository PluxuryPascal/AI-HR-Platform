package domain

type CandidateCreatedEvent struct {
	CandidateID   string `json:"candidate_id"`
	JobID         string `json:"job_id"`
	ResumeFileKey string `json:"resume_file_key"`
}
