package cmd

import (
	"errors"
	"fmt"
)

var ErrCreateLogFileFailed = errors.New("create log file failed")
var ErrLogFileDisabled = errors.New("logfile disabled")

type logErr struct {
	Key   error
	Cause error
}

func wrapLogErr(key error, cause error) error {
	return &logErr{
		Key:   key,
		Cause: cause,
	}
}

func (e *logErr) Is(target error) bool {
	return e.Key == target
}

func (e *logErr) Unwrap() error {
	return e.Cause
}

func (e *logErr) Error() string {
	var prefix string
	var suffix string
	if e.Key != nil {
		prefix = fmt.Sprintf("%s: ", e.Key.Error())
	}
	if e.Cause != nil {
		suffix = e.Cause.Error()
	}
	return prefix + suffix
}

func (e *logErr) String() string {
	return e.Error()
}
