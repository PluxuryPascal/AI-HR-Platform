package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CommunicationRepository interface {
	Create(ctx context.Context, c *domain.Communication) error
	GetByID(ctx context.Context, id string) (*domain.Communication, error)
	GetByCandidateID(ctx context.Context, candidateID string) ([]domain.Communication, error)
	MarkSent(ctx context.Context, communicationID string) error
}

type communicationRepo struct {
	dbClient *db.PostgresClient
}

func NewCommunicationRepo(dbClient *db.PostgresClient) CommunicationRepository {
	return &communicationRepo{dbClient: dbClient}
}

func (r *communicationRepo) Create(ctx context.Context, c *domain.Communication) error {
	const query = `
		INSERT INTO ai_engine.t_communications (candidate_id, generated_by_user_id, type, subject, body)
		VALUES (@candidate_id, @generated_by_user_id, @type, @subject, @body)
		RETURNING id, created_at
	`

	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"candidate_id":       c.CandidateID,
		"generated_by_user_id": c.GeneratedByUserID,
		"type":               c.Type,
		"subject":            c.Subject,
		"body":               c.Body,
	}).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert communication: %w", err)
	}

	return nil
}

func (r *communicationRepo) GetByID(ctx context.Context, id string) (*domain.Communication, error) {
	const query = `
		SELECT id, candidate_id, generated_by_user_id, type, subject, body, sent_at, created_at
		FROM ai_engine.t_communications
		WHERE id = @id
	`

	var c domain.Communication
	err := r.dbClient.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id": id,
	}).Scan(
		&c.ID, &c.CandidateID, &c.GeneratedByUserID, &c.Type,
		&c.Subject, &c.Body, &c.SentAt, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("communication not found")
		}
		return nil, fmt.Errorf("scan communication: %w", err)
	}

	return &c, nil
}

func (r *communicationRepo) GetByCandidateID(ctx context.Context, candidateID string) ([]domain.Communication, error) {
	const query = `
		SELECT id, candidate_id, generated_by_user_id, type, subject, body, sent_at, created_at
		FROM ai_engine.t_communications
		WHERE candidate_id = @candidate_id
		ORDER BY created_at DESC
	`

	rows, err := r.dbClient.Pool.Query(ctx, query, pgx.NamedArgs{
		"candidate_id": candidateID,
	})
	if err != nil {
		return nil, fmt.Errorf("query communications: %w", err)
	}
	defer rows.Close()

	var comms []domain.Communication
	for rows.Next() {
		var c domain.Communication
		if err := rows.Scan(
			&c.ID, &c.CandidateID, &c.GeneratedByUserID, &c.Type,
			&c.Subject, &c.Body, &c.SentAt, &c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan communication: %w", err)
		}
		comms = append(comms, c)
	}

	return comms, nil
}

func (r *communicationRepo) MarkSent(ctx context.Context, communicationID string) error {
	const query = `
		UPDATE ai_engine.t_communications
		SET sent_at = NOW()
		WHERE id = @id AND sent_at IS NULL
	`

	_, err := r.dbClient.Pool.Exec(ctx, query, pgx.NamedArgs{
		"id": communicationID,
	})
	if err != nil {
		return fmt.Errorf("mark communication sent: %w", err)
	}

	return nil
}
