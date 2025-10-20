package env

import (
	"os"
	"strings"
)

var MvnCommandTimeout = envi("MVN_COMMAND_TIMEOUT", 0)
var DisableMvnCommand = strings.TrimSpace(os.Getenv("NO_MVN")) != ""
var MavenCentral string

func init() {
	if strings.TrimSpace(os.Getenv("SKIP_MAVEN_CENTRAL")) == "" {
		MavenCentral = "https://repo1.maven.org/maven2/"
	}
}
