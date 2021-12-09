package version

import (
	"fmt"
	"github.com/MakeNowJust/heredoc/v2"
)

const version = "v"

func Version() string {
	return version
}

func PrintVersionInfo() {
	fmt.Printf(heredoc.Doc(`
	murphysec-cli %s
	`), version)
}
