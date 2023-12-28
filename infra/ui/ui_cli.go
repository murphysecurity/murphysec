package ui

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"os"
)

type cli struct{}

var CLI UI = &cli{}

var _ UI = (*cli)(nil)
var Term = termenv.NewOutput(os.Stdout)

func (cli) UpdateStatus(s Status, msg string) {
	cliStatus = s
	cliStatusMsg = msg
	if isatty.IsTerminal(os.Stdout.Fd()) {
		Term.ClearLine()
	}
	fmt.Print("\r")
	statusRepaint()
}

func (cli) Display(level MessageLevel, msg string) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		Term.ClearLine()
	}
	if level == MsgError {
		fmt.Println(Term.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s] %s", level.String(), msg)))
	} else {
		fmt.Println(Term.String().Foreground(level.fColor()).Styled(fmt.Sprintf("[%s]", level.String())), msg)
	}
	statusRepaint()
}

func (cli) ClearStatus() {
	if cliStatus == StatusIdle {
		return
	}
	cliStatus = StatusIdle
	cliStatusMsg = ""
	if isatty.IsTerminal(os.Stdout.Fd()) {
		Term.ClearLine()
	}
}

var cliStatus = StatusIdle
var cliStatusMsg = ""

func statusRepaint() {
	if cliStatus == StatusIdle {
		return
	}
	fmt.Print(Term.String().Foreground(cliStatus.fColor()).Styled(cliStatus.String()))
	if cliStatusMsg != "" {
		fmt.Print(" - ", cliStatusMsg)
	}
	fmt.Print("\r")
}
