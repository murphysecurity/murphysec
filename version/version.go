package version

import (
	"fmt"
	"github.com/iseki0/osname"
	"github.com/murphysecurity/murphysec/utils/must"
	"os"
	"path/filepath"
)

const version = "v1.8.2"

// PrintVersionInfo print version info to stdout
func PrintVersionInfo() {
	fmt.Printf("%s %s\n", filepath.Base(must.A(filepath.EvalSymlinks(must.A(os.Executable())))), Version())
}

var userAgent string

func init() {
	osn, e := osname.OsName()
	if e != nil {
		osn = "<unknownOS>"
	}
	userAgent = fmt.Sprintf("murphysec-cli/%s (%s);", Version(), osn)
}

func UserAgent() string {
	return userAgent
}
