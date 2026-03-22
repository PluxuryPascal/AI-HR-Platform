package domain

type CandidateCreatedEvent struct {
	CandidateID   string `json:"candidate_id"`
	JobID         string `json:"job_id"`
	ResumeFileKey string `json:"resume_file_key"`
	TeamID        string `json:"team_id"`
	Locale        string `json:"locale"`
}

type InviteCreatedEvent struct {
	InviteID string `json:"invite_id"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	Role     string `json:"role"`
	Locale   string `json:"locale"`
}
