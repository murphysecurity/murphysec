//go:build windows

package display

import "github.com/muesli/termenv"

func EnableANSI() {
	termenv.EnableWindowsANSIConsole()
}
