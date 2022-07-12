package env

import (
	"os"
	"strconv"
	"strings"
)

var GradleExecutionTimeoutSecond = envi("GRADLE_EXECUTION_TIMEOUT_SEC", 20*60)

func envi(name string, defaultValue int) int {
	if i, e := strconv.Atoi(os.Getenv(name)); e != nil {
		return defaultValue
	} else {
		return i
	}
}

var SpecificProjectName = ""
var DisableGit = false
var _ServerBaseURL = "https://www.murphysec.com"

func ConfigureServerBaseUrl(u string) {
	_ServerBaseURL = strings.TrimRight(strings.TrimSpace(u), "/")
}

func ServerBaseUrl() string {
	return _ServerBaseURL
}
