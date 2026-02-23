package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultEncoding   = "json"
	defaultLevel      = zapcore.InfoLevel
	defaultOutputPath = "stdout"
	defaultWriteStd   = false
)

type Log struct {
	Log *zap.Logger

	files []*os.File

	level   zapcore.Level
	cores   []zapcore.Core
	encoder zapcore.Encoder
}

func New(opts ...Option) (*Log, error) {
	l := &Log{
		level: defaultLevel,
		cores: make([]zapcore.Core, 0),
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	l.encoder = zapcore.NewJSONEncoder(encoderConfig)

	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	if len(l.cores) == 0 {
		return nil, fmt.Errorf("no cores configured")
	}

	core := zapcore.NewTee(l.cores...)
	l.Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return l, nil
}
