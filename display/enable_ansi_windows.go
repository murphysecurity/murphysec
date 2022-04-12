//go:build windows

package display

import "github.com/muesli/termenv"

var _ = func() int {
	termenv.EnableWindowsANSIConsole()
	return 0
}()
