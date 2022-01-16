package logger

import (
	"log"
	"murphysec-cli-simple/utils/must"
	"os"
	"strings"
	"sync"
)

type LogLevel int

const (
	LogDebug LogLevel = iota + 1
	LogInfo
	LogWarn
	LogErr
	LogSilent
)

var logMutex sync.Mutex
var consoleLogLevel = LogWarn

var ConsoleLogLevelOverride string

var getConsoleLogLevel = func() func() LogLevel {
	o := sync.Once{}
	c := LogWarn
	return func() LogLevel {
		o.Do(func() {
			switch strings.ToLower(strings.TrimSpace(ConsoleLogLevelOverride)) {
			case "error":
				c = LogErr
			case "warn":
				c = LogWarn
			case "info":
				c = LogInfo
			case "debug":
				c = LogDebug
			case "silent":
				c = LogSilent
			case "":
			default:
				panic("loglevel invalid")
			}
		})
		return c
	}
}()

var Debug = log.New(n(LogDebug), "[DEBUG]", log.Lshortfile+log.LstdFlags)
var Info = log.New(n(LogInfo), "[INFO]", log.Lshortfile+log.LstdFlags)
var Warn = log.New(n(LogWarn), "[WARN]", log.Lshortfile+log.LstdFlags)
var Err = log.New(n(LogErr), "[ERROR]", log.Lshortfile+log.LstdFlags)

type W struct {
	l LogLevel
}

func n(l LogLevel) *W {
	return &W{l: l}
}

func (w *W) Write(p []byte) (n int, err error) {
	logMutex.Lock()
	defer logMutex.Unlock()
	if w.l >= getConsoleLogLevel() {
		must.Int(os.Stderr.Write(p))
	}
	f := loggerFile()
	if f != nil {
		must.Int(f.Write(p))
	}
	return len(p), nil
}
