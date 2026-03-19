package domain

import "time"

type TeamAISettings struct {
	TeamID     string    `json:"team_id"`
	APIKey     *string   `json:"api_key,omitempty"`
	ParseModel *string   `json:"parse_model,omitempty"`
	ScoreModel *string   `json:"score_model,omitempty"`
	EmbedModel *string   `json:"embed_model,omitempty"`
	ChatModel  *string   `json:"chat_model,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
