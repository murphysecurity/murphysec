package ui

import (
	"fmt"
	"github.com/muesli/termenv"
)

type CLI struct{}

var _ UI = (*CLI)(nil)

func (CLI) UpdateStatus(s Status, msg string) {
	cliStatus = s
	cliStatusMsg = msg
	termenv.ClearLine()
	fmt.Print("\r")
	statusRepaint()
}

func (CLI) Display(level MessageLevel, msg string) {
	termenv.ClearLine()
	if level == MsgError {
		fmt.Println(termenv.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s] %s", level.String(), msg)))
	} else {
		fmt.Println(termenv.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s]", level.String())), msg)
	}
	statusRepaint()
}

func (CLI) ClearStatus() {
	cliStatus = StatusIdle
	cliStatusMsg = ""
	termenv.ClearLine()
}

var cliStatus = StatusIdle
var cliStatusMsg = ""

func statusRepaint() {
	if cliStatus == StatusIdle {
		return
	}
	fmt.Print(termenv.String().Foreground(cliStatus.fColor()).Styled(cliStatus.String()))
	if cliStatusMsg != "" {
		fmt.Print(" - ", cliStatusMsg)
	}
	fmt.Print("\r")
}
