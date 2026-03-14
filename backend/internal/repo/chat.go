package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ChatRepository interface {
	SaveSession(ctx context.Context, session *domain.ChatSession) error
	AddMessage(ctx context.Context, message *domain.ChatMessage) error
	GetHistory(ctx context.Context, sessionID string) ([]domain.ChatMessage, error)
}

type chatRepo struct {
	dbClient *db.PostgresClient
}

func NewChatRepo(dbClient *db.PostgresClient) ChatRepository {
	return &chatRepo{dbClient: dbClient}
}

func (r *chatRepo) SaveSession(ctx context.Context, session *domain.ChatSession) error {
	const query = `
		INSERT INTO ai_engine.t_chat_sessions (team_id, user_id, type, target_candidate_id, title)
		VALUES (@team_id, @user_id, @type, @target, @title)
		RETURNING id, created_at, updated_at
	`
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"team_id":  session.TeamID,
		"user_id":  session.UserID,
		"type":     session.Type,
		"target":   session.TargetCandidateID,
		"title":    session.Title,
	}).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return fmt.Errorf("save chat session: %w", err)
	}
	return nil
}

func (r *chatRepo) AddMessage(ctx context.Context, message *domain.ChatMessage) error {
	const query = `
		INSERT INTO ai_engine.t_chat_messages (session_id, role, content, tokens_used)
		VALUES (@session_id, @role, @content, @tokens)
		RETURNING id, created_at
	`
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"session_id": message.SessionID,
		"role":       message.Role,
		"content":    message.Content,
		"tokens":     message.TokensUsed,
	}).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return fmt.Errorf("add chat message: %w", err)
	}
	return nil
}

func (r *chatRepo) GetHistory(ctx context.Context, sessionID string) ([]domain.ChatMessage, error) {
	const query = `
		SELECT id, session_id, role, content, tokens_used, created_at
		FROM ai_engine.t_chat_messages
		WHERE session_id = @session_id
		ORDER BY created_at ASC
	`
	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("query chat history: %w", err)
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var m domain.ChatMessage
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.TokensUsed, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan chat message: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, nil
}
