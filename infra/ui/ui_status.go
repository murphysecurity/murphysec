package ui

import "github.com/muesli/termenv"

//go:generate stringer -type Status -linecomment -output ui_status_string.go
type Status int

const (
	_               Status = iota
	StatusIdle             // IDLE
	StatusRunning          // RUNNING
	StatusWaiting          // WAITING
	StatusSucceeded        // SUCCEEDED
	StatusFailed           // FAILED
)

func (i Status) fColor() termenv.ANSIColor {
	switch i {
	case StatusFailed:
		return termenv.ANSIBrightRed
	case StatusSucceeded:
		return termenv.ANSIBrightGreen
	case StatusWaiting, StatusRunning:
		return termenv.ANSIBrightBlue
	default:
		return termenv.ANSIBrightCyan
	}
}
