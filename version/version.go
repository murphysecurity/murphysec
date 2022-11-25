package version

import (
	"github.com/iseki0/osname"
	"github.com/murphysecurity/murphysec/infra/buildinfo"
	"runtime"
	"strings"
)

const name = "murphysec-cli"

var userAgent string

func initUserAgent() {
	userAgent += name + "/" + Version()
	var (
		e         error
		platforms []string
	)
	osn, e := osname.OsName()
	if e != nil && osn != "" {
		osn = "<unknownOS>"
	}
	platforms = append(platforms, osn, runtime.GOOS, runtime.GOARCH)
	userAgent += " " + strings.Join(platforms, ";")
	if s := buildinfo.UserAgentSuffix(); s != "" {
		userAgent += " " + s
	}
}

func init() {
	initUserAgent()
}

func UserAgent() string {
	return userAgent
}

func Version() string {
	return buildinfo.Get().Version // + "-" + build_flags.Kind
}

func FullInfo() string {
	return name + " " + Version() + "\n\nBuild: " + buildinfo.Commit()
}
