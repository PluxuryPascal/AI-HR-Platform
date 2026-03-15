package temporal

import (
	"backend/internal/temporal/activity"
	"backend/internal/temporal/workflow"
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"

	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

type Worker struct {
	log            *zap.Logger
	workerCfg      *worker.Options
	cfg            *config.Temporal
	temporalClient *Client
	TemporalWorker worker.Worker
	activities     *activity.Activities
}

func NewWorker(log *zap.Logger, temporalClient *Client, cfg *config.Temporal, activities *activity.Activities) *Worker {
	return &Worker{
		log: log,
		workerCfg: &worker.Options{
			DeploymentOptions: worker.DeploymentOptions{},
		},
		temporalClient: temporalClient,
		activities:     activities,
		cfg:            cfg,
	}
}

func (w *Worker) DependsOn() []string {
	return []string{"temporal-client", "logger", "db", "cloudinary"}
}

func (w *Worker) HealthCheck(ctx context.Context) error {
	return nil
}

func (w *Worker) Init(ctx context.Context) error {
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize:     w.cfg.WorkerCount,
		MaxConcurrentWorkflowTaskExecutionSize: 100,
	}

	w.TemporalWorker = worker.New(w.temporalClient.TemporalClient, w.cfg.QueueName, workerOptions)

	w.TemporalWorker.RegisterWorkflow(workflow.ResumePipelineWorkflow)
	w.TemporalWorker.RegisterActivity(w.activities)

	return nil
}

func (w *Worker) Name() string {
	return "temporal-worker"
}

func (w *Worker) Run(ctx context.Context) error {
	if err := w.TemporalWorker.Start(); err != nil {
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
