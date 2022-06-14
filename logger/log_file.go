package logger

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/utils/must"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var CliLogFilePathOverride string
var DisableLogFile bool

var defaultLogFile = filepath.Join(must.A(homedir.Dir()), ".murphysec", "logs", fmt.Sprintf("%d.log", time.Now().UnixMilli()))

func CreateLogFile() (*os.File, error) {
	if DisableLogFile {
		return nil, ErrLogFileDisabled
	}
	logFilePath := defaultLogFile
	if CliLogFilePathOverride != "" {
		logFilePath = CliLogFilePathOverride
	}
	// ensure log dir created
	if e := os.MkdirAll(filepath.Dir(logFilePath), 0755); e != nil {
		return nil, &LogErr{
			Key:   ErrCreateLogFileFailed,
			Cause: e,
		}
	}
	if f, e := os.OpenFile(logFilePath, os.O_CREATE+os.O_RDWR+os.O_APPEND, 0644); e != nil {
		return nil, &LogErr{
			Key:   ErrCreateLogFileFailed,
			Cause: e,
		}
	} else {
		return f, nil
	}
}

// file before staticRefTime will be ignored
var staticRefTime = must.A(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z"))

// LogFileCleanup auto remove log files which created between staticRefTime and 7 days ago
func LogFileCleanup() {
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
			if lt.Before(staticRefTime) {
				continue
			}
			if time.Now().Sub(time.UnixMilli(int64(ts))) > time.Hour*24*7 {
				_ = os.Remove(filepath.Join(basePath, entry.Name()))
			}
		}
	}
}
