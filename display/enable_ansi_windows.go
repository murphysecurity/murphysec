//go:build windows

package display

import "github.com/muesli/termenv"

func init() {
	_, _ = termenv.EnableWindowsANSIConsole()
}
