package logger

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var CliLogFilePathOverride string
var DisableLogFile bool

var loggerFile = func() func() *os.File {
	o := sync.Once{}
	var file *os.File
	f := func() *os.File {
		o.Do(func() {
			if DisableLogFile {
				return
			}
			var logFilePath = filepath.Join(must.String(homedir.Dir()), ".murphysec", "logs", fmt.Sprintf("%d.log", time.Now().UnixMilli()))
			if CliLogFilePathOverride != "" {
				logFilePath = CliLogFilePathOverride
			}
			if e := os.MkdirAll(filepath.Dir(logFilePath), 0755); e != nil {
				panic(errors.Wrap(e, fmt.Sprintf("Create log file dirs failed. %s", e.Error())))
			}
			_f, e := os.OpenFile(logFilePath, os.O_CREATE+os.O_RDWR+os.O_APPEND, 0644)
			if e != nil {
				panic(errors.Wrap(e, fmt.Sprintf("Open log file failed. %s, %s", logFilePath, e.Error())))
			}
			file = _f
		})
		return file
	}
	return f
}()
