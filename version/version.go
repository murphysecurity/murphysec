package version

import (
	"fmt"
	"murphysec-cli-simple/util/must"
	"os"
	"path/filepath"
)

const version = "1.1.3"

// Version returns version string
func Version() string {
	return version
}

// PrintVersionInfo print version info to stdout
func PrintVersionInfo() {
	fmt.Printf("%s %s\n", filepath.Base(must.String(filepath.EvalSymlinks(must.String(os.Executable())))), version)
}
