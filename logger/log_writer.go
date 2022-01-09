package logger

import (
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"sync"
)

var fileMutex = sync.Mutex{}
var consoleMutex = sync.Mutex{}
var fdv *os.File
var fileLogLevel = LogDebug
var consoleLogLevel = LogWarn

func SetConsoleLogLevel(level LogLevel) {
	consoleLogLevel = level
}
func InitLogFile(filename string) {
	must.Must(os.MkdirAll(filepath.Dir(filename), 0755))
	if fdv != nil {
		panic(1)
	}
	f, e := os.OpenFile(filename, os.O_CREATE+os.O_APPEND, 0644)
	must.Must(e)
	fdv = f
}

type LogWriter struct {
	lv LogLevel
}

func (this *LogWriter) Write(b []byte) (int, error) {
	if this.lv >= consoleLogLevel {
		func() {
			consoleMutex.Lock()
			defer consoleMutex.Unlock()
			_, e := os.Stderr.Write(b)
			must.Must(e)
		}()
	}
	if fdv != nil && this.lv >= fileLogLevel {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, _ = fdv.Write(b)
	}
	return len(b), nil
}
func NewLogWriter(lv LogLevel) *LogWriter {
	return &LogWriter{lv: lv}
}
