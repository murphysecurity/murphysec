package utils

import (
	"context"
	"go.uber.org/zap"
)

const _LoggerCtxKey = `_LoggerCtxKey`

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, _LoggerCtxKey, logger)
}

// UseLogger returns *zap.Logger in the context. If no Logger exists, returns zap.NewNop()
func UseLogger(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(_LoggerCtxKey).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}
