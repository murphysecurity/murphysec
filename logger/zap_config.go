package logger

import "go.uber.org/zap/zapcore"

var ZapConsoleLoggerEncoderConfig = zapcore.EncoderConfig{
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
	EncodeName:          zapcore.FullNameEncoder,
	NewReflectedEncoder: nil,
	ConsoleSeparator:    " ",
}

var ZapConsoleEncoder = zapcore.NewConsoleEncoder(ZapConsoleLoggerEncoderConfig)
