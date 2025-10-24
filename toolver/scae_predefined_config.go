package toolver

import (
	"context"
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

func locateJavaHome(version string) string {
	var javaHome string
	switch version {
	case "jdk8":
		javaHome = "/opt/java/8"
	case "jdk11":
		javaHome = "/opt/java/11"
	case "jdk17":
		javaHome = "/opt/java/17"
	case "jdk21":
		javaHome = "/opt/java/21"
	}
	if javaHome != "" {
		return javaHome
	}
	return ""
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
