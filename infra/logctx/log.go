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
	return ctx.Value(key).(*zap.Logger)
}
