package display

import (
	"fmt"
	"github.com/muesli/termenv"
)

type _CLI struct{}

func (_ _CLI) UpdateStatus(s Status, msg string) {
	cliStatus = s
	cliStatusMsg = msg
	// todo
	termenv.ClearLine()
	fmt.Print("\r")
	statusRepaint()
}

func (_ _CLI) Display(level MsgLevel, msg string) {
	termenv.ClearLine()
	if level == MsgError {
		fmt.Println(termenv.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s] %s", level.String(), msg)))
	} else {
		fmt.Println(termenv.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s]", level.String())), msg)
	}
	statusRepaint()
}

func (_ _CLI) ClearStatus() {
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
