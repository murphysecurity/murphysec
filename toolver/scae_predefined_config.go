package toolver

import (
	"context"
	"os"
)

func init() {
	postprocessors = append(postprocessors, func(ctx context.Context, config *Config) {
		if config.Maven.JdkVersion != "" && config.Maven.JavaHome == "" {
			config.Maven.JavaHome = locateJavaHome(config.Maven.JdkVersion)
		}
		if config.Maven.MavenVersion != "" && config.Maven.MavenCommand == "" {
			config.Maven.MavenCommand = locateMvnHome(config.Maven.MavenVersion) + "/bin/mvn"
		}
	})
}

var javaHomeMap = map[string]string{
	"jdk8":  "/opt/java/8",
	"jdk11": "/opt/java/11",
	"jdk17": "/opt/java/17",
	"jdk21": "/opt/java/21",
	"jdk25": "/opt/java/25",
}

func init() {
	for version, home := range javaHomeMap {
		_, e := os.Stat(home)
		if os.IsNotExist(e) {
			delete(javaHomeMap, version)
		}
	}
}

func locateJavaHome(version string) string {
	return javaHomeMap[version]
}

func locateMvnHome(version string) string {
	var mvnHome string
	switch version {
	case "maven3.6.3":
		mvnHome = "/opt/maven/3.6.3"
	case "maven3.8.8":
		mvnHome = "/opt/maven/3.8.8"
	case "maven3.9.5":
		mvnHome = "/opt/maven/3.9.5"
	}
	if mvnHome != "" {
		return mvnHome
	}
	return ""
}
