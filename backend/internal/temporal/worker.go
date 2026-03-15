package temporal

import (
	"backend/internal/temporal/activity"
	"backend/pkg/svc"
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

type Worker struct {
	log            *zap.Logger
	workerCfg      *worker.Options
	queueName      string
	temporalClient client.Client
	TemporalWorker worker.Worker
	activities     activity.Activities
}

func NewWorker(log *zap.Logger, temporalClient client.Client, queueName string, activities activity.Activities) *Worker {
	return &Worker{
		log: log,
		workerCfg: &worker.Options{
			DeploymentOptions: worker.DeploymentOptions{},
		},
		temporalClient: temporalClient,
		activities:     activities,
		queueName:      queueName,
	}
}

func (w *Worker) DependsOn() []string {
	return []string{"temporal-client", "logger"}
}

func (w *Worker) HealthCheck(ctx context.Context) error {
	return nil
}

func (w *Worker) Init(ctx context.Context) error {
	w.TemporalWorker = worker.New(w.temporalClient, w.queueName, *w.workerCfg)

	return nil
}

func (w *Worker) Name() string {
	return "temporal-worker"
}

func (w *Worker) Run(ctx context.Context) error {
	if err := w.TemporalWorker.Run(worker.InterruptCh()); err != nil {
		return fmt.Errorf("temporal worker failed: %w", err)
	}

	return nil
}

func (w *Worker) Stop(ctx context.Context) error {
	if w.TemporalWorker != nil {
		w.TemporalWorker.Stop()
	}

	return nil
}

var _ svc.Service = (*Worker)(nil)
