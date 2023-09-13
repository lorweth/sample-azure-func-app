package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Infof(msg string, args ...interface{})

	Errorf(err error, msg string, args ...interface{})

	With(fields ...zap.Field) Logger
}
