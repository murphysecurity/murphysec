//go:build windows

package ui

import "github.com/muesli/termenv"

func init() {
	_, _ = termenv.EnableWindowsANSIConsole()
}
