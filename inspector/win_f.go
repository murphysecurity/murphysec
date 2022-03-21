//go:build windows

package inspector

import "github.com/muesli/termenv"

func EnableANSI() {
	termenv.EnableWindowsANSIConsole()
}
