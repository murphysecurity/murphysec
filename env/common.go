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

var DisableMvnCommand = strings.TrimSpace(os.Getenv("NO_MVN")) != ""
var IdeaInstallPath = os.Getenv("IDEA_INSTALL")
var MvnCommandTimeout = envi("MVN_COMMAND_TIMEOUT", 0)
var MavenCentral string

func init() {
	if strings.TrimSpace(os.Getenv("SKIP_MAVEN_CENTRAL")) == "" {
		MavenCentral = "https://repo1.maven.org/maven2/"
	}
}

func ConfigureServerBaseUrl(u string) {
	_ServerBaseURL = strings.TrimRight(strings.TrimSpace(u), "/")
}

func ServerBaseUrl() string {
	return _ServerBaseURL
}
