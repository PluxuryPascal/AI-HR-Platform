package notification

import (
	"backend/internal/domain"
	"backend/pkg/email"
	"backend/pkg/mq"
	"context"
	"encoding/json"
	"fmt"

	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

// InviteCreatedHandler возвращает mq.HandlerFunc для обработки событий invite.created.
//
// При получении сообщения:
//  1. Десериализует InviteCreatedEvent из JSON.
//  2. Формирует текст приглашения.
//  3. Отправляет email через emailClient.Send().
//
// Ошибки парсинга → NackDiscard (повторные попытки бессмысленны).
// Ошибки отправки → NackRequeue (временная ошибка, retry имеет смысл).
func InviteCreatedHandler(emailClient *email.Client, log *zap.Logger) mq.HandlerFunc {
	return func(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action {
		var event domain.InviteCreatedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Error("failed to unmarshal invite created event",
				zap.Error(err),
				zap.ByteString("body", d.Body),
			)

			return rabbitmq.NackDiscard
		}

		log.Info("processing invite created event",
			zap.String("invite_id", event.InviteID),
			zap.String("email", event.Email),
			zap.String("role", event.Role),
		)

		subject := "You've been invited to join AI-HR Platform"
		body := fmt.Sprintf(
			"Hello!\n\nYou have been invited to join AI-HR Platform as %s.\n\nUse the following token to complete your registration: %s\n\nThis invitation will expire soon, so please act promptly.\n\nBest regards,\nAI-HR Team",
			event.Role,
			event.Token,
		)

		if err := emailClient.Send(ctx, event.Email, subject, body); err != nil {
			log.Error("failed to send invite email",
				zap.Error(err),
				zap.String("email", event.Email),
				zap.String("invite_id", event.InviteID),
			)

			return rabbitmq.NackRequeue
		}

		log.Info("invite email sent successfully",
			zap.String("email", event.Email),
			zap.String("invite_id", event.InviteID),
		)

		return rabbitmq.Ack
	}
}
