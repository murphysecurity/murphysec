package version

import (
	"fmt"
	"github.com/iseki0/osname"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/utils/must"
	"os"
	"path/filepath"
	"sync"
)

const version = "v1.6.3"

// PrintVersionInfo print version info to stdout
func PrintVersionInfo() {
	fmt.Printf("%s %s\n", filepath.Base(must.A(filepath.EvalSymlinks(must.A(os.Executable())))), Version())
}

var _ua = func() func() string {
	o := sync.Once{}
	ua := ""
	osn, e := osname.OsName()
	if e != nil {
		osn = "<unknownOS>"
	}
	return func() string {
		o.Do(func() {
			ua = fmt.Sprintf("murphysec-cli/%s (%s);", Version(), osn)
			logger.Debug.Println("user-agent:", ua)
		})
		return ua
	}
}()

func UserAgent() string {
	return _ua()
}
