package version

import (
	"fmt"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"sync"
)

const version = "1.2.0-alpha"

// Version returns version string
func Version() string {
	return version
}

// PrintVersionInfo print version info to stdout
func PrintVersionInfo() {
	fmt.Printf("%s %s\n", filepath.Base(must.String(filepath.EvalSymlinks(must.String(os.Executable())))), version)
}

var _ua = func() func() string {
	o := sync.Once{}
	ua := ""
	return func() string {
		o.Do(func() {
			ua = fmt.Sprintf("murphysec-cli/%s (%s);", Version(), getOSVersion())
			logger.Debug.Println("user-agent:", ua)
		})
		return ua
	}
}()

func UserAgent() string {
	return _ua()
}
