package domain

import (
	"time"
)

// ChatType represents the chat_type ENUM
type ChatType string

const (
	ChatLocalCandidate ChatType = "local_candidate"
	ChatGlobalSearch   ChatType = "global_search"
)

// MessageRole represents the message_role ENUM
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// ChatSession represents the ai_engine.t_chat_sessions table
type ChatSession struct {
	ID                string     `json:"id" db:"id"`
	TeamID            string     `json:"team_id" db:"team_id"`
	UserID            string     `json:"user_id" db:"user_id"`
	Type              ChatType   `json:"type" db:"type"`
	TargetCandidateID *string    `json:"target_candidate_id,omitempty" db:"target_candidate_id"`
	Title             *string    `json:"title,omitempty" db:"title"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// ChatMessage represents the ai_engine.t_chat_messages table
type ChatMessage struct {
	ID         string      `json:"id" db:"id"`
	SessionID  string      `json:"session_id" db:"session_id"`
	Role       MessageRole `json:"role" db:"role"`
	Content    string      `json:"content" db:"content"`
	TokensUsed *int        `json:"tokens_used,omitempty" db:"tokens_used"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
}
