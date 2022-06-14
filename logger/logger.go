package logger

import (
	"fmt"
	"github.com/pkg/errors"
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

var loggerEncoderConfig = zapcore.EncoderConfig{
	MessageKey:          "message",
	LevelKey:            "level",
	TimeKey:             "time",
	NameKey:             "name",
	CallerKey:           "caller",
	FunctionKey:         "",
	StacktraceKey:       "stacktrace",
	SkipLineEnding:      false,
	LineEnding:          zapcore.DefaultLineEnding,
	EncodeLevel:         zapcore.CapitalLevelEncoder,
	EncodeTime:          zapcore.RFC3339TimeEncoder,
	EncodeDuration:      zapcore.StringDurationEncoder,
	EncodeCaller:        zapcore.ShortCallerEncoder,
	EncodeName:          nil,
	NewReflectedEncoder: nil,
	ConsoleSeparator:    " ",
}

func InitLogger() error {
	// console logger
	encoder := zapcore.NewConsoleEncoder(loggerEncoderConfig)
	consoleCore := zapcore.NewNopCore()
	switch strings.ToLower(strings.TrimSpace(ConsoleLogLevelOverride)) {
	case "error":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), zapcore.ErrorLevel)
	case "warn":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), zapcore.WarnLevel)
	case "info":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), zapcore.InfoLevel)
	case "debug":
		consoleCore = zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), zapcore.DebugLevel)
	}

	// file logger
	var fileCore zapcore.Core
	logFile, e := CreateLogFile()
	if e == nil {
		fileCore = zapcore.NewCore(encoder, logFile, zapcore.DebugLevel)
	} else if errors.Is(e, ErrLogFileDisabled) {
		fileCore = zapcore.NewNopCore()
	} else {
		return e
	}

	// all
	core := zapcore.NewTee(fileCore, consoleCore)

	logger := zap.New(core, zap.AddCaller())
	wl := logger.WithOptions(zap.AddCallerSkip(1))
	Debug = w(logger, wl.Debug)
	Info = w(logger, wl.Info)
	Warn = w(logger, wl.Warn)
	Err = w(logger, wl.Warn)
	Logger = logger.Sugar()

	NetworkLogger = zap.NewNop().Sugar()

	return nil
}

var Debug = noOp
var Info = noOp
var Warn = noOp
var Err = noOp
