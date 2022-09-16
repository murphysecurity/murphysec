package version

import (
	"fmt"
	"github.com/iseki0/osname"
	"runtime"
	"runtime/debug"
)

const version = "v1.10.0"

var userAgent string

func init() {
	osn, e := osname.OsName()
	if e != nil {
		osn = "<unknownOS>"
	}
	var platform = fmt.Sprintf("%s; %s; %s", osn, runtime.GOOS, runtime.GOARCH)
	userAgent = fmt.Sprintf("murphysec-cli/%s (%s)", Version(), platform)
	if h := GetGitHash(); h != "" {
		userAgent = userAgent + " GitHash/" + h
	}
	if h := GetGitTime(); h != "" {
		userAgent = userAgent + " GitTime/" + GetGitTime()
	}
}

func UserAgent() string {
	return userAgent
}

var buildInfo map[string]string

func fillBuildInfo() {
	if buildInfo != nil {
		return
	}
	info, b := debug.ReadBuildInfo()
	buildInfo = map[string]string{}
	if !b {
		return
	}
	for _, it := range info.Settings {
		buildInfo[it.Key] = it.Value
	}
}

func GetGitHash() string {
	fillBuildInfo()
	if buildInfo["vcs"] != "git" {
		return ""
	}
	return buildInfo["vcs.revision"]
}

func GetGitModified() string {
	return buildInfo["vcs.modified"]
}

func GetGitTime() string {
	return buildInfo["vcs.time"]
}
