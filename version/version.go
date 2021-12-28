package version

import (
	"fmt"
	"github.com/MakeNowJust/heredoc/v2"
)

const version = "1.1.9"

// Version returns version string
func Version() string {
	return version
}

// PrintVersionInfo print version info to stdout
func PrintVersionInfo() {
	fmt.Printf(heredoc.Doc(`
	murphysec %s
	`), version)
}
