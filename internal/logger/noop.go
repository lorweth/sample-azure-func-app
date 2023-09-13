package logger

import (
	"go.uber.org/zap"
)

type noopLogger struct{}

func (n noopLogger) Infof(msg string, args ...interface{}) {}

func (n noopLogger) Errorf(err error, msg string, args ...interface{}) {}

func (n noopLogger) With(fields ...zap.Field) Logger {
	return noopLogger{}
}
