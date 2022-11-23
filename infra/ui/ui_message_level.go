package ui

import "github.com/muesli/termenv"

//go:generate stringer -type MessageLevel -linecomment -output ui_message_level_string.go
type MessageLevel int

const (
	_         MessageLevel = iota
	MsgInfo                // Info
	MsgNotice              // Notice
	MsgWarn                // Warn
	MsgError               // Error
)

func (i MessageLevel) fColor() termenv.ANSIColor {
	switch i {
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
