package mq

import "go.uber.org/zap"

// zapLogger адаптирует *zap.SugaredLogger под интерфейс rabbitmq.Logger.
// Интерфейс требует: Fatalf / Errorf / Warnf / Infof / Debugf —
// все они есть у SugaredLogger «из коробки», поэтому адаптер тривиален.
type zapLogger struct {
	s *zap.SugaredLogger
}

func newZapLogger(l *zap.Logger) *zapLogger {
	return &zapLogger{s: l.Sugar()}
}

func (z *zapLogger) Debugf(format string, args ...interface{}) { z.s.Debugf(format, args...) }
func (z *zapLogger) Infof(format string, args ...interface{})  { z.s.Infof(format, args...) }
func (z *zapLogger) Warnf(format string, args ...interface{})  { z.s.Warnf(format, args...) }
func (z *zapLogger) Errorf(format string, args ...interface{}) { z.s.Errorf(format, args...) }
func (z *zapLogger) Fatalf(format string, args ...interface{}) { z.s.Fatalf(format, args...) }
