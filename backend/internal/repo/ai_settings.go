package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type AiSettingsRepository interface {
	GetByTeamID(ctx context.Context, teamID string) (*domain.TeamAISettings, error)
	Upsert(ctx context.Context, settings *domain.TeamAISettings) error
}

type aiSettingsRepo struct {
	dbClient *db.PostgresClient
}

func NewAiSettingsRepo(dbClient *db.PostgresClient) AiSettingsRepository {
	return &aiSettingsRepo{dbClient: dbClient}
}

func (r *aiSettingsRepo) GetByTeamID(ctx context.Context, teamID string) (*domain.TeamAISettings, error) {
	const query = `
		SELECT team_id, api_key, parse_model, score_model, embed_model, chat_model, created_at, updated_at
		FROM ai_engine.t_team_settings
		WHERE team_id = @team_id
	`
	var s domain.TeamAISettings
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{"team_id": teamID}).Scan(
		&s.TeamID, &s.APIKey, &s.ParseModel, &s.ScoreModel, &s.EmbedModel, &s.ChatModel, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Return nil if settings not found
		}
		return nil, fmt.Errorf("query ai settings: %w", err)
	}

	return &s, nil
}

func (r *aiSettingsRepo) Upsert(ctx context.Context, settings *domain.TeamAISettings) error {
	const query = `
		INSERT INTO ai_engine.t_team_settings (team_id, api_key, parse_model, score_model, embed_model, chat_model, created_at, updated_at)
		VALUES (@team_id, @api_key, @parse_model, @score_model, @embed_model, @chat_model, NOW(), NOW())
		ON CONFLICT (team_id) DO UPDATE SET
			api_key = EXCLUDED.api_key,
			parse_model = EXCLUDED.parse_model,
			score_model = EXCLUDED.score_model,
			embed_model = EXCLUDED.embed_model,
			chat_model = EXCLUDED.chat_model,
			updated_at = NOW()
		RETURNING created_at, updated_at
	`
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"team_id":     settings.TeamID,
		"api_key":     settings.APIKey,
		"parse_model": settings.ParseModel,
		"score_model": settings.ScoreModel,
		"embed_model": settings.EmbedModel,
		"chat_model":  settings.ChatModel,
	}).Scan(&settings.CreatedAt, &settings.UpdatedAt)

	if err != nil {
		return fmt.Errorf("upsert ai settings: %w", err)
	}

	return nil
}
