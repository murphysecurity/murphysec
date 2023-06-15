package ui

import (
	"testing"
)

func TestName(t *testing.T) {
	cli{}.Display(MsgError, "Error message")
	cli{}.Display(MsgWarn, "Warn message")
	cli{}.Display(MsgNotice, "Notice message")
	cli{}.Display(MsgInfo, "Info message")
}
