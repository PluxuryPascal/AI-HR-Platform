package logger

import (
	"backend/pkg/svc"
	"context"
	"fmt"
)

var _ svc.Service = (*Log)(nil)

const errSync = "sync /dev/stdout: invalid argument; sync /dev/stderr: invalid argument"

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
	if err := l.Log.Sync(); err != nil && err.Error() != errSync {
		return fmt.Errorf("failed to sync logger: %w", err)
	}
	return nil
}
