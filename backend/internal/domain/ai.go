package domain

type ParseResult struct {
	FirstName       *string  `json:"first_name,omitempty"`
	LastName        *string  `json:"last_name,omitempty"`
	Email           *string  `json:"email,omitempty"`
	Location        *string  `json:"location,omitempty"`
	Skills          []string `json:"skills"`
	ExperienceYears *int     `json:"experience_years,omitempty"`
	Education       *string  `json:"education,omitempty"`
	Summary         *string  `json:"summary,omitempty"`
	StructuredData  []byte   `json:"structured_data"`
	JobRequirements string   `json:"job_requirements"`
	MissingFields   []string `json:"missing_fields"`
}

func (r *ParseResult) ParsingStatus() CandidateParsingStatus {
	for _, f := range r.MissingFields {
		if f == "email" || f == "first_name" || f == "last_name" || f == "location" {
			return ParsingStatusNeedsReview
		}
	}

	return ParsingStatusCompleted
}

type ScoreResult struct {
	MatchScore int           `json:"match_score"` // 0–100
	Factors    []ScoreFactor `json:"factors"`
}

type EmbeddingChunk struct {
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding"`
}

type JobParseResult struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Requirements []string `json:"requirements"`
	WorkFormat   string   `json:"work_format"`
	SalaryMin    int    `json:"salary_min"`
	SalaryMax    int    `json:"salary_max"`
	Currency     string `json:"currency"`
}
