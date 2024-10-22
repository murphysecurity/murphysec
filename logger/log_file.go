package logger

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/utils/must"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)

const defaultLogFilePath = env.HomeSubdir + "/logs"

// CreateLogFile create log file. If _filepath is empty, use default log path
func CreateLogFile(_filepath string) (_ *os.File, err error) {
	defer func() {
		if err != nil {
			err = errors.WithCause(ErrCreateLogFileFailed, err)
		}
	}()
	var logFilepath = _filepath
	if logFilepath == "" {
		if home, e := homedir.Dir(); e != nil {
			return nil, e
		} else {
			logFilepath = filepath.Join(home, defaultLogFilePath, fmt.Sprintf("%d.log", time.Now().UnixMilli()))
		}
	}
	if e := os.MkdirAll(filepath.Dir(logFilepath), 0755); e != nil {
		return nil, e
	}
	_ = checkOrCreateCacheDirTag(filepath.Dir(logFilepath))
	if f, e := os.OpenFile(logFilepath, os.O_CREATE+os.O_RDWR+os.O_APPEND, 0644); e != nil {
		return nil, e
	} else {
		return f, nil
	}
}

// LogFileCleanup auto remove log files which created between staticRefTime and 7 days ago
func LogFileCleanup() {
	// file before staticRefTime will be ignored
	var staticRefTime = must.A(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z"))

	logFilePattern := regexp.MustCompile(`^(\d+)\.log$`)
	home, e := homedir.Dir()
	if e != nil {
		return
	}
	basePath := filepath.Join(home, defaultLogFilePath)
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
			if time.Since(time.UnixMilli(int64(ts))) > time.Hour*24*7 {
				_ = os.Remove(filepath.Join(basePath, entry.Name()))
			}
		}
	}
}

const cacheDirTagPrefix = `Signature: 8a477f597d28d172789f06886806bc55`
const cacheDirTagSuffix = `
# This file is a cache directory tag created by murphysec.
# For information about cache directory tags, see:
#	https://bford.info/cachedir/
`
const cacheDirTagFileName = "CACHEDIR.TAG"

func checkOrCreateCacheDirTag(dir string) (e error) {
	var binaryCDT = unsafe.Slice(unsafe.StringData(cacheDirTagPrefix), len(cacheDirTagPrefix))
	var f *os.File
	var ccFp = filepath.Join(dir, cacheDirTagFileName)
	f, e = os.Open(ccFp)
	var fileOk = false
	if e == nil {
		var buf = make([]byte, len(cacheDirTagPrefix))
		_, e = io.ReadFull(f, buf)
		if e == nil && bytes.Compare(binaryCDT, buf) == 0 {
			fileOk = true
		}
		e = f.Close()
		if e != nil {
			return
		}
	}
	if fileOk {
		return
	}
	var binaryCDT2 = unsafe.Slice(unsafe.StringData(cacheDirTagSuffix), len(cacheDirTagSuffix))
	f, e = os.OpenFile(ccFp, os.O_CREATE+os.O_WRONLY+os.O_TRUNC, 0644)
	if e != nil {
		return
	}
	defer func() {
		var e2 = f.Close()
		if e == nil {
			e = e2
		}
	}()
	_, e = f.Write(binaryCDT)
	if e != nil {
		return
	}
	_, e = f.Write(binaryCDT2)
	return
}
