package gradle

import (
	"fmt"
	"murphysec-cli-simple/util/output"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// detectGradleVersion returns gradle version
func detectGradleVersion(path string) (*gradleVersion, error) {
	c := exec.Command(path, "-v")
	output.Debug(fmt.Sprintf("Run: %s", c.String()))
	rs, e := c.Output()
	if e != nil {
		output.Info(fmt.Sprintf("detectGradleVersion failed, %v", e))
		return nil, e
	}
	output.Debug("gradle version output:")
	output.Debug(string(rs))
	v := parseGradleVersionOutput(string(rs))
	return &v, nil
}

func parseGradleVersionOutput(s string) gradleVersion {
	rs := gradleVersion{}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "Gradle"):
			rs.Version = strings.TrimSpace(strings.Trim(line, "Gradle"))
		case strings.HasPrefix(line, "Build time:"):
			rs.BuildTime = strings.TrimSpace(strings.Trim(line, "Build time:"))
		case strings.HasPrefix(line, "Revision:"):
			rs.Revision = strings.TrimSpace(strings.Trim(line, "Revision:"))
		case strings.HasPrefix(line, "Kotlin:"):
			rs.Kotlin = strings.TrimSpace(strings.Trim(line, "Kotlin: "))
		case strings.HasPrefix(line, "Groovy:"):
			rs.Groovy = strings.TrimSpace(strings.Trim(line, "Groovy:"))
		case strings.HasPrefix(line, "Ant:"):
			rs.Ant = strings.TrimSpace(strings.Trim(line, "Ant:"))
		case strings.HasPrefix(line, "JVM:"):
			rs.JVM = strings.TrimSpace(strings.Trim(line, "JVM:"))
		case strings.HasPrefix(line, "OS:"):
			rs.OS = strings.TrimSpace(strings.Trim(line, "OS:"))
		}
	}
	return rs
}

type gradleVersion struct {
	Version   string
	BuildTime string
	Revision  string
	Kotlin    string
	Groovy    string
	Ant       string
	JVM       string
	OS        string
}

// detectGradleFile returns gradle file path in dir, returns empty if not found.
func detectGradleFile(dir string) string {
	for s := range gradleFiles {
		p := filepath.Join(dir, s)
		output.Debug(fmt.Sprintf("try to detect gradle file: %s", p))
		if stat, err := os.Stat(filepath.Join(dir, s)); err == nil && !stat.IsDir() {
			output.Debug("found")
			return p
		}
	}
	output.Debug(fmt.Sprintf("not found any gradle file under: %s", dir))
	return ""
}

// detectGradleWrapper returns gradle wrapper script path if a suitable gradle wrapper script exists
func detectGradleWrapper(dir string) string {
	f := ""
	if runtime.GOOS == "windows" {
		f = "gradlew.bat"
	} else {
		f = "gradlew"
	}
	p := filepath.Join(dir, f)
	if s, e := os.Stat(p); e == nil && !s.IsDir() {
		return p
	} else {
		return ""
	}
}

func getGradleCmd(dir string) string {
	if p := detectGradleWrapper(dir); p != "" {
		return p
	}
	return "gradle" // cli fallback
}

var gradleFiles = map[string]bool{
	"build.gradle":     true,
	"build.gradle.kts": true,
}

var skipDetectDirs = []string{
	".git",
	"node_modules",
	".idea",
	".gradle",
}
