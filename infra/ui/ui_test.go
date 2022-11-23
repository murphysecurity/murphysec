package ui

import (
	"testing"
)

func TestName(t *testing.T) {
	CLI{}.Display(MsgError, "Error message")
	CLI{}.Display(MsgWarn, "Warn message")
	CLI{}.Display(MsgNotice, "Notice message")
	CLI{}.Display(MsgInfo, "Info message")
}
