package logctx

import (
	"context"
	"go.uber.org/zap"
)

type _key struct{}

var key = &_key{}

func With(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

func Use(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(key).(*zap.Logger)
	if !ok || l == nil {
		return zap.NewNop()
	}
	return l
}
