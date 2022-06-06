package logger

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var CliLogFilePathOverride string
var DisableLogFile bool

var defaultLogFile = filepath.Join(must.A(homedir.Dir()), ".murphysec", "logs", fmt.Sprintf("%d.log", time.Now().UnixMilli()))
var loggerFile = func() func() *os.File {
	o := sync.Once{}
	var file *os.File
	f := func() *os.File {
		o.Do(func() {
			if DisableLogFile {
				return
			}
			var logFilePath = defaultLogFile
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

func LogFileCleanup() {
	refTime, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	logFilePattern := regexp.MustCompile("^(\\d+)\\.log$")
	basePath := filepath.Dir(defaultLogFile)
	if basePath == "" {
		return
	}
	d, e := os.ReadDir(basePath)
	if e != nil {
		return
	}
	for _, entry := range d {
		if entry.IsDir() || !entry.Type().IsRegular() {
			continue
		}
		if m := logFilePattern.FindStringSubmatch(entry.Name()); m != nil {
			ts, e := strconv.Atoi(m[1])
			if e != nil {
				continue
			}
			lt := time.UnixMilli(int64(ts))
			if lt.Before(refTime) {
				continue
			}

			if time.Now().Sub(time.UnixMilli(int64(ts))) > time.Hour*24*7 {
				_ = os.Remove(filepath.Join(basePath, entry.Name()))
			}
		}
	}
}
