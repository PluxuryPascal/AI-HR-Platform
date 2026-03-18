package worker

import (
	"backend/internal/domain"
	"backend/internal/temporal"
	"backend/internal/temporal/activity"
	"backend/internal/temporal/workflow"
	"context"
	"encoding/json"
	"fmt"
	"time"

	hiringv1 "backend/internal/proto/hiring/v1"

	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

func NewCandidateConsumerHandler(log *zap.Logger, temporalClient *temporal.Client) func(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action {
	return func(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action {
		var event domain.CandidateCreatedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Error("failed to unmarshal candidate created event", zap.Error(err))
			return rabbitmq.NackDiscard
		}

		log.Info("received candidate.created event",
			zap.String("candidate_id", event.CandidateID),
			zap.String("job_id", event.JobID),
			zap.String("team_id", event.TeamID),
		)

		input := activity.ResumePipelineInput{
			CandidateID:   event.CandidateID,
			JobID:         event.JobID,
			ResumeFileKey: event.ResumeFileKey,
			TeamID:        event.TeamID,
			Locale:        event.Locale,
		}

		workflowID := fmt.Sprintf("resume-pipeline-%s", event.CandidateID)

		_, err := temporalClient.StartWorkflow(ctx, workflowID, 30*time.Minute, workflow.ResumePipelineWorkflow, input)
		if err != nil {
			log.Error("failed to start resume pipeline workflow",
				zap.String("candidate_id", event.CandidateID),
				zap.Error(err),
			)
			return rabbitmq.NackRequeue
		}

		log.Info("resume pipeline workflow started",
			zap.String("candidate_id", event.CandidateID),
			zap.String("workflow_id", workflowID),
		)

		return rabbitmq.Ack
	}
}

func NewDLQHandler(log *zap.Logger, temporalClient *temporal.Client) func(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action {
	return func(ctx context.Context, d rabbitmq.Delivery) rabbitmq.Action {
		var event domain.CandidateCreatedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Error("failed to unmarshal dead letter event", zap.Error(err))
			return rabbitmq.NackDiscard
		}

		log.Warn("candidate processing failed permanently, sending FAILED callback",
			zap.String("candidate_id", event.CandidateID),
			zap.String("job_id", event.JobID),
		)

		input := activity.GRPCCallbackInput{
			CandidateID: event.CandidateID,
			JobID:       event.JobID,
			Status:      hiringv1.ParsingStatus_PARSING_STATUS_FAILED,
		}

		workflowID := fmt.Sprintf("dlq-callback-%s", event.CandidateID)

		_, err := temporalClient.StartWorkflow(ctx, workflowID, 5*time.Minute, workflow.DLQCallbackWorkflow, input)
		if err != nil {
			log.Error("failed to start DLQ callback workflow",
				zap.String("candidate_id", event.CandidateID),
				zap.Error(err),
			)
			return rabbitmq.NackRequeue
		}

		return rabbitmq.Ack
	}
}
