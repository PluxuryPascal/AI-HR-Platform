package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*Log) error

func WithLevel(level string) Option {
	return func(l *Log) error {
		lvl, err := extractLevel(level)
		if err != nil {
			return fmt.Errorf("failed to extract level: %w", err)
		}

		l.level = lvl

		return nil
	}
}

func WithStdOut(enabled bool) Option {
	return func(l *Log) error {
		if !enabled {
			return nil
		}

		coreStdout := zapcore.NewCore(
			l.encoder,
			zapcore.AddSync(os.Stdout),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.DebugLevel && lvl <= zapcore.InfoLevel
			}),
		)

		coreStdErr := zapcore.NewCore(
			l.encoder,
			zapcore.AddSync(os.Stderr),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.WarnLevel
			}),
		)

		l.cores = append(l.cores, coreStdout, coreStdErr)

		return nil
	}
}

func WithFile(path string) Option {
	return func(l *Log) error {
		path, err := checkOrCreatePath(path)
		if err != nil {
			return fmt.Errorf("failed to check or create path: %w", err)
		}

		cleanPath := filepath.Clean(path)

		file, err := os.OpenFile(cleanPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		core := zapcore.NewCore(
			l.encoder,
			zapcore.AddSync(file),
			l.level,
		)

		l.cores = append(l.cores, core)

		return nil
	}
}
