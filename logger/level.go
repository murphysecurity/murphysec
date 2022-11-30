package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

//go:generate stringer -type Level -output level_string.go -linecomment
type Level int

const (
	_           Level = iota
	LevelSilent       // silent
	LevelDebug        // debug
	LevelInfo         // info
	LevelWarn         // warn
	LevelError        // error
)

func (i Level) Valid() bool {
	return i > 0 && int(i) < len(_Level_index)
}

func (i *Level) Of(s string) error {
	switch strings.ToLower(s) {
	case "silent", "":
		*i = LevelSilent
	case "debug":
		*i = LevelDebug
	case "info":
		*i = LevelInfo
	case "warn":
		*i = LevelWarn
	case "error":
		*i = LevelError
	default:
		return errors.New("bad loglevel")
	}
	return nil
}

func (i Level) ZapLevel() zapcore.Level {
	switch i {
	case LevelError:
		return zap.ErrorLevel
	case LevelWarn:
		return zap.WarnLevel
	case LevelInfo:
		return zap.InfoLevel
	case LevelDebug:
		return zap.DebugLevel
	}
	return zap.DebugLevel
}
