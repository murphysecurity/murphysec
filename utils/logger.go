package utils

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return logctx.With(ctx, logger)
}

// UseLogger returns *zap.Logger in the context. If no Logger exists, returns zap.NewNop()
func UseLogger(ctx context.Context) *zap.Logger {
	return logctx.Use(ctx)
}
