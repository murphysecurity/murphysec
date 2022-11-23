package display

import (
	"github.com/muesli/termenv"
	"strings"
)

var (
	CLI  UI = _CLI{}
	NONE UI = _NONE{}
)

type Status int

const (
	StatusIdle Status = iota
	StatusRunning
	StatusWaiting
	StatusSucceeded
	StatusFailed
)

func (s Status) String() string {
	switch s {
	case StatusFailed:
		return "FAILED"
	case StatusIdle:
		return "IDLE"
	case StatusRunning:
		return "RUNNING"
	case StatusSucceeded:
		return "SUCCEEDED"
	case StatusWaiting:
		return "WAITING"
	}
	panic("")
}

func (s Status) fColor() termenv.ANSIColor {
	switch s {
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

type MsgLevel int

func (m *MsgLevel) UnmarshalText(text []byte) error {
	t := strings.ToLower(string(text))
	switch t {
	case "info":
		*m = MsgInfo
	case "notice":
		*m = MsgNotice
	case "warn":
		*m = MsgWarn
	case "error":
		*m = MsgError
	default:
		*m = MsgNotice
	}
	return nil
}

const (
	MsgInfo MsgLevel = iota
	MsgNotice
	MsgWarn
	MsgError
)

func (m MsgLevel) fColor() termenv.ANSIColor {
	switch m {
	case MsgInfo:
		return termenv.ANSIBrightGreen
	case MsgError:
		return termenv.ANSIBrightRed
	case MsgWarn:
		return termenv.ANSIBrightRed
	case MsgNotice:
		return termenv.ANSIBrightCyan
	}
	panic(0)
}

func (m MsgLevel) String() string {
	switch m {
	case MsgInfo:
		return "Info"
	case MsgWarn:
		return "Warn"
	case MsgNotice:
		return "Notice"
	case MsgError:
		return "Error"
	}
	panic("")
}

type UI interface {
	UpdateStatus(s Status, msg string)
	Display(level MsgLevel, msg string)
	ClearStatus()
}
