package version

import (
	"fmt"
	"github.com/iseki0/osname"
)

const version = "v1.9.7"

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
