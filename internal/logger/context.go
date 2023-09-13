package logger

import (
	"context"
)

const (
	loggerCtxKey = "logger"
)

func FromCtx(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerCtxKey).(Logger)
	if !ok {
		return noopLogger{}
	}

	return l
}

func SetInCtx(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

func NewCtx(ctx context.Context) context.Context {
	return SetInCtx(context.Background(), FromCtx(ctx))
}
