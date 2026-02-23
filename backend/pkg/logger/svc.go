package logger

import (
	"backend/pkg/svc"
	"context"
	"errors"
	"fmt"
	"syscall"
)

var _ svc.Service = (*Log)(nil)

func (l *Log) DependsOn() []string {
	return nil
}

func (l *Log) HealthCheck(ctx context.Context) error {
	return nil
}

func (l *Log) Name() string {
	return "logger"
}

func (l *Log) Init(ctx context.Context) error {
	return nil
}

func (l *Log) Run(ctx context.Context) error {
	return nil
}

func (l *Log) Stop(ctx context.Context) error {
	err := l.Log.Sync()

	if err != nil && !errors.Is(err, syscall.EINVAL) {
		return fmt.Errorf("failed to sync logger: %w", err)
	}

	return nil
}
