package logger

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type WrappedLogger struct {
	*zap.Logger
	f func(msg string, fields ...zap.Field)
}

var noOp = &WrappedLogger{
	Logger: zap.NewNop(),
	f:      func(msg string, fields ...zap.Field) {},
}

var Debug = noOp
var Info = noOp
var Warn = noOp
var Err = noOp

func (this *WrappedLogger) Printf(format string, a ...interface{}) {
	this.f(fmt.Sprintf(strings.TrimSuffix(format, "\n"), a...))
}

func (this *WrappedLogger) Println(args ...interface{}) {
	this.f(fmt.Sprint(args...))
}

func w(logger *zap.Logger, f func(msg string, fields ...zap.Field)) *WrappedLogger {
	return &WrappedLogger{
		Logger: logger,
		f:      f,
	}
}

func InitLegacyLogger(logger *zap.Logger) {
	wl := logger.WithOptions(zap.AddCallerSkip(1))
	Debug = w(logger, wl.Debug)
	Info = w(logger, wl.Info)
	Warn = w(logger, wl.Warn)
	Err = w(logger, wl.Warn)
}
