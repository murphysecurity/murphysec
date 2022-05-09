package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var ConsoleLogLevelOverride string
var Logger *zap.SugaredLogger
var NetworkLogger *zap.SugaredLogger

type WrappedLogger struct {
	*zap.Logger
	f func(msg string, fields ...zap.Field)
}

var noOp = &WrappedLogger{
	Logger: zap.NewNop(),
	f:      func(msg string, fields ...zap.Field) {},
}

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

func InitLogger() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "name",
		CallerKey:           "caller",
		FunctionKey:         "",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          "",
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.RFC3339TimeEncoder,
		EncodeDuration:      zapcore.StringDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		EncodeName:          nil,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    " ",
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileCore := zapcore.NewCore(encoder, zapcore.Lock(loggerFile()), zapcore.DebugLevel)
	consoleCore := zapcore.NewNopCore()
	switch strings.ToLower(strings.TrimSpace(ConsoleLogLevelOverride)) {
	case "silent":
		consoleCore = zapcore.NewNopCore()
	case "error":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.ErrorLevel)
	case "warn":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.WarnLevel)
	case "info":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)
	case "debug":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	}
	core := zapcore.NewTee(fileCore, consoleCore)

	logger := zap.New(core, zap.AddCaller())
	wl := logger.WithOptions(zap.AddCallerSkip(1))
	Debug = w(logger, wl.Debug)
	Info = w(logger, wl.Info)
	Warn = w(logger, wl.Warn)
	Err = w(logger, wl.Warn)
	Logger = logger.Sugar()

	NetworkLogger = zap.NewNop().Sugar()
}

var Debug = noOp
var Info = noOp
var Warn = noOp
var Err = noOp
