package ui

import (
	"fmt"
	"github.com/muesli/termenv"
)

type cli struct{}

var CLI UI = &cli{}

var _ UI = (*cli)(nil)

func (cli) UpdateStatus(s Status, msg string) {
	cliStatus = s
	cliStatusMsg = msg
	termenv.ClearLine() //nolint:all
	fmt.Print("\r")
	statusRepaint()
}

func (cli) Display(level MessageLevel, msg string) {
	termenv.ClearLine() //nolint:all
	if level == MsgError {
		fmt.Println(termenv.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s] %s", level.String(), msg)))
	} else {
		fmt.Println(termenv.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s]", level.String())), msg)
	}
	statusRepaint()
}

func (cli) ClearStatus() {
	if cliStatus == StatusIdle {
		return
	}
	cliStatus = StatusIdle
	cliStatusMsg = ""
	termenv.ClearLine() //nolint:all
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
