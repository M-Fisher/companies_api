package logger

import (
	"context"

	"go.uber.org/zap"
)

type loggerContenttype int

const loggerKey loggerContenttype = iota

var defLogger *zap.Logger

func SetDefaultLogger(l *zap.Logger) {
	defLogger = l
}

func GetDefaultLogger() *zap.Logger {
	if defLogger != nil {
		return defLogger
	}
	return zap.L()
}

func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return l
	}
	return GetDefaultLogger()
}
