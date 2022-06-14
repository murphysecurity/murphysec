package logger

import (
	"errors"
	"fmt"
)

var ErrCreateLogFileFailed = errors.New("create log file failed")
var ErrLogFileDisabled = errors.New("logfile disabled")

type LogErr struct {
	Key   error
	Cause error
}

func (e *LogErr) Is(target error) bool {
	return e.Key == target
}

func (e *LogErr) Unwrap() error {
	return e.Cause
}

func (e *LogErr) Error() string {
	var prefix string
	var suffix string
	if e.Key != nil {
		prefix = fmt.Sprintf("%s: ", e.Key.Error())
	}
	if e.Cause != nil {
		suffix = e.Error()
	}
	return prefix + suffix
}

func (e *LogErr) String() string {
	return e.Error()
}
