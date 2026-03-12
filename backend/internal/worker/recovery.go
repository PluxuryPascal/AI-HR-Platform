package worker

import (
	"backend/internal/usecase"
	"backend/pkg/config"
	"backend/pkg/svc"
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type InviteRecoveryWorker struct {
	log     *zap.Logger
	cfg     *config.InviteRecovery
	usecase usecase.InviteUseCase
	cron    *cron.Cron
}

func NewInviteRecoveryWorker(log *zap.Logger, cfg *config.InviteRecovery, usecase usecase.InviteUseCase) *InviteRecoveryWorker {
	return &InviteRecoveryWorker{
		log:     log,
		cfg:     cfg,
		usecase: usecase,
		cron:    cron.New(),
	}
}

func (w *InviteRecoveryWorker) Name() string {
	return "invite-recovery-worker"
}

func (w *InviteRecoveryWorker) DependsOn() []string {
	return []string{"logger", "db", "api"} // Зависит от API, так как API создает начальные записи
}

func (w *InviteRecoveryWorker) Init(ctx context.Context) error {
	_, err := w.cron.AddFunc(w.cfg.Cron, func() {
		w.log.Debug("running invite recovery cycle")
		if err := w.usecase.ProcessStuckInvites(context.Background()); err != nil {
			w.log.Error("failed to process stuck invites", zap.Error(err))
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	return nil
}

func (w *InviteRecoveryWorker) HealthCheck(ctx context.Context) error {
	return nil
}

func (w *InviteRecoveryWorker) Run(ctx context.Context) error {
	w.cron.Start()

	return nil
}

func (w *InviteRecoveryWorker) Stop(ctx context.Context) error {
	w.cron.Stop()

	return nil
}

var _ svc.Service = (*InviteRecoveryWorker)(nil)
