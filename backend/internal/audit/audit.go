package audit

import (
	"backend/internal/db"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// ActorType mirrors the actor_type ENUM in the database.
type ActorType string

const (
	ActorUser    ActorType = "user"
	ActorAI      ActorType = "ai_agent"
	ActorSystem  ActorType = "system"
)

// Entry represents a single audit log entry to be recorded.
type Entry struct {
	TeamID    string
	ActorType ActorType
	ActorID   *string   // nil for system/ai actors
	Action    string    // action code, e.g. "hiring.job_created"
	TargetID  *string   // optional target entity UUID
	Payload   any       // optional, will be marshalled to JSONB
}

// Logger provides audit logging via direct writes to logging.t_activity_logs.
type Logger struct {
	log      *zap.Logger
	dbClient *db.PostgresClient
	// actionCache maps action code → action_id for fast lookups.
	actionCache map[string]int
}

// NewLogger creates a new audit Logger.
func NewLogger(log *zap.Logger, dbClient *db.PostgresClient) *Logger {
	return &Logger{
		log:         log,
		dbClient:    dbClient,
		actionCache: make(map[string]int),
	}
}

// SeedActionTypes idempotently inserts all known action types and populates the cache.
func (l *Logger) SeedActionTypes(ctx context.Context) error {
	const upsert = `
		INSERT INTO logging.t_action_types (service, code, description)
		VALUES (@service, @code, @description)
		ON CONFLICT (code) DO UPDATE SET
			service = EXCLUDED.service,
			description = EXCLUDED.description
		RETURNING id
	`

	for _, a := range AllActions() {
		var id int
		err := l.dbClient.Pool.QueryRow(ctx, upsert, pgx.NamedArgs{
			"service":     a.Service,
			"code":        a.Code,
			"description": a.Description,
		}).Scan(&id)
		if err != nil {
			return fmt.Errorf("seed action %q: %w", a.Code, err)
		}
		l.actionCache[a.Code] = id
	}

	l.log.Info("audit action types seeded", zap.Int("count", len(l.actionCache)))
	return nil
}

// Log writes an audit entry to logging.t_activity_logs.
// It resolves the action code to an ID via the internal cache.
// If the action code is unknown, it logs a warning and returns nil.
func (l *Logger) Log(ctx context.Context, e Entry) error {
	actionID, ok := l.actionCache[e.Action]
	if !ok {
		l.log.Warn("unknown audit action code, skipping", zap.String("action", e.Action))
		return nil
	}

	var payloadBytes []byte
	if e.Payload != nil {
		var err error
		payloadBytes, err = json.Marshal(e.Payload)
		if err != nil {
			l.log.Warn("failed to marshal audit payload", zap.Error(err))
			// Don't fail the business operation because of audit serialization.
		}
	}

	// Determine service from the action code prefix.
	service := serviceFromAction(e.Action)

	const query = `
		INSERT INTO logging.t_activity_logs
			(team_id, service, actor_type, actor_id, action_id, target_id, payload)
		VALUES
			(@team_id, @service, @actor_type, @actor_id, @action_id, @target_id, @payload)
	`

	_, err := l.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"team_id":    e.TeamID,
		"service":    service,
		"actor_type": e.ActorType,
		"actor_id":   e.ActorID,
		"action_id":  actionID,
		"target_id":  e.TargetID,
		"payload":    payloadBytes,
	})
	if err != nil {
		l.log.Warn("failed to write audit log",
			zap.String("action", e.Action),
			zap.Error(err),
		)
		// Return nil to avoid failing the business operation.
		return nil
	}

	return nil
}

// serviceFromAction extracts the service name from an action code.
// e.g. "hiring.job_created" → "hiring", "ai.resume_parsed" → "ai_engine".
func serviceFromAction(code string) string {
	for i, r := range code {
		if r == '.' {
			prefix := code[:i]
			if prefix == "ai" {
				return "ai_engine"
			}
			return prefix
		}
	}
	return "unknown"
}
