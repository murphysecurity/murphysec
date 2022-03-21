package display

import (
	"fmt"
	"github.com/muesli/termenv"
	"os"
)

var (
	CLI  UI = _CLI{}
	NONE UI = _NONE{}
)

type _CLI struct{}
type _NONE struct{}

func (_ _NONE) ClearStatus() {}

func (_ _NONE) UpdateStatus(s Status, msg string) {}

func (_ _NONE) Display(level MsgLevel, msg string) {}

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

const (
	MsgInfo MsgLevel = iota
	MsgNotice
	MsgWarn
	MsgError
)

func (m MsgLevel) fColor() termenv.ANSIColor {
	switch m {
	case MsgInfo:
		return termenv.ANSIBrightWhite
	case MsgError:
		return termenv.ANSIBrightRed
	case MsgWarn:
		return termenv.ANSIBrightYellow
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

var cliStatus = StatusIdle
var cliStatusMsg = ""

func (_ _CLI) ClearStatus() {
	cliStatus = StatusIdle
	cliStatusMsg = ""
}

func statusRepaint() {
	if cliStatus == StatusIdle {
		return
	}
	// todo: disable it for some old terminal
	if os.Getenv("TERM_PROGRAM") == "Apple_Terminal" {
		fmt.Println(cliStatus.String(), cliStatusMsg)
		return
	}
	termenv.SaveCursorPosition()
	fmt.Print(termenv.String().Foreground(cliStatus.fColor()).Styled(cliStatus.String()))
	if cliStatusMsg != "" {
		fmt.Print(" - ", cliStatusMsg)
	}
	termenv.RestoreCursorPosition()
}
func (_ _CLI) UpdateStatus(s Status, msg string) {
	cliStatus = s
	cliStatusMsg = msg
	// todo
	if os.Getenv("TERM_PROGRAM") == "Apple_Terminal" {
		statusRepaint()
		return
	}
	termenv.ClearLine()
	statusRepaint()
}

func (_ _CLI) Display(level MsgLevel, msg string) {
	if os.Getenv("TERM_PROGRAM") == "Apple_Terminal" {
		fmt.Println(level.String(), msg)
		return
	}
	termenv.ClearLine()
	fmt.Println(termenv.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s]", level.String())), msg, "\r")
	statusRepaint()
}
