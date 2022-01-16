package logger

import (
	"log"
)

type LogLevel int

const (
	LogDebug LogLevel = iota + 1
	LogInfo
	LogWarn
	LogErr
)

var debugP = NewLogWriter(LogDebug)
var infoP = NewLogWriter(LogInfo)
var warnP = NewLogWriter(LogWarn)
var errP = NewLogWriter(LogErr)
var Debug = log.New(debugP, "[DEBUG]", log.Lshortfile+log.LstdFlags)
var Info = log.New(infoP, "[INFO]", log.Lshortfile+log.LstdFlags)
var Warn = log.New(warnP, "[WARN]", log.Lshortfile+log.LstdFlags)
var Err = log.New(errP, "[ERROR]", log.Lshortfile+log.LstdFlags)
