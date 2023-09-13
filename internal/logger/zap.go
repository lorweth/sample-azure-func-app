package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type structuredLogger struct {
	zapLogger *zap.Logger
	logTags   map[string]string
}

func New() Logger {
	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(
		enc,
		os.Stdout,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
	)

	return structuredLogger{
		zapLogger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel)),
	}
}

func (s structuredLogger) Infof(format string, args ...interface{}) {
	s.zapLogger.Info(fmt.Sprintf(format, args...))
}

func (s structuredLogger) Errorf(err error, format string, args ...interface{}) {
	// Append otel exception keys
	zapLogger := s.zapLogger.With(
		zap.String("exception.type", fmt.Sprintf("%T", err)),
		zap.String("exception.message", err.Error()),
	)

	zapLogger.Error(fmt.Sprintf(format+" %+v", append(args, err)...))
}

func (s structuredLogger) With(fields ...zap.Field) Logger {
	return &structuredLogger{
		zapLogger: s.zapLogger.With(fields...),
	}
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}
