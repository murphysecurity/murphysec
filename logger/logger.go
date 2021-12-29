package logger

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
	"log"
	"murphysec-cli-simple/util/must"
	"os"
	"path/filepath"
	"time"
)

var Debug = log.New(os.Stderr, "[DEBUG]", log.LstdFlags)
var Info = log.New(os.Stderr, "[INFO]", log.LstdFlags)
var Warn = log.New(os.Stderr, "[WARN]", log.LstdFlags)
var Err = log.New(os.Stderr, "[ERROR]", log.LstdFlags)

func InitFileLog(path string) error {
	if path == "" {
		path = filepath.Join(must.String(homedir.Dir()), ".murphysec", "logs", fmt.Sprintf("%d.log", time.Now().UnixMilli()))
	}
	f, e := os.OpenFile(path, os.O_CREATE+os.O_APPEND, 0644)
	if e != nil {
		return e
	}
	Debug.SetOutput(NewCopyWriter(os.Stderr, f))
	Info.SetOutput(NewCopyWriter(os.Stderr, f))
	Warn.SetOutput(NewCopyWriter(os.Stderr, f))
	Err.SetOutput(NewCopyWriter(os.Stderr, f))
	return nil
}

type CopyWriter struct {
	writer []io.Writer
}

func (this CopyWriter) Write(p []byte) (n int, err error) {
	for _, it := range this.writer {
		_, e := it.Write(p)
		if e != nil {
			return 0, e
		}
	}
	return len(p), nil
}
func NewCopyWriter(w ...io.Writer) *CopyWriter {
	return &CopyWriter{w}
}
