package domain

import (
	"time"
)

type ActorType string

const (
	ActorTypeUser    ActorType = "user"
	ActorTypeAIAgent ActorType = "ai_agent"
	ActorTypeSystem  ActorType = "system"
)

// ActionType represents the logging.t_action_types table
type ActionType struct {
	ID          int     `json:"id" db:"id"`
	Service     string  `json:"service" db:"service"` // 'auth', 'hiring', 'ai_engine'
	Code        string  `json:"code" db:"code"`       // Unique code, e.g., 'hiring.candidate_moved'
	Description *string `json:"description,omitempty" db:"description"`
}

// ActivityLog represents the logging.t_activity_logs table
type ActivityLog struct {
	ID        string    `json:"id" db:"id"`
	TeamID    string    `json:"team_id" db:"team_id"`
	Service   string    `json:"service" db:"service"`
	ActorType ActorType `json:"actor_type" db:"actor_type"`
	ActorID   *string   `json:"actor_id,omitempty" db:"actor_id"`
	ActionID  int       `json:"action_id" db:"action_id"`
	TargetID  *string   `json:"target_id,omitempty" db:"target_id"`
	Payload   []byte    `json:"payload,omitempty" db:"payload"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
