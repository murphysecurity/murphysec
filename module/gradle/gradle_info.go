package gradle

import (
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

//goland:noinspection GoNameStartsWithPackageName
type GradleInfo struct {
	Executable string
	Version    string
	Revision   string
	CallCmd    func(args ...string) *exec.Cmd
}

func (g *GradleInfo) String() string {
	return fmt.Sprintf("Gradle[%s]: %s , revision: %s", g.Version, g.Executable, g.Revision)
}

func evalGradleInfo(dir string) (*GradleInfo, error) {
	info, e := execWrappedGradleInfo(dir)
	if e != nil {
		logger.Debug.Println("check gradle wrapper failed.", e.Error())
	} else {
		return info, nil
	}
	info, e = execRawGradleInfo(dir)
	if e != nil {
		logger.Debug.Println("check raw gradle failed.", e.Error())
	} else {
		return info, nil
	}
	return nil, e
}

func parseGradleVersion(s string) GradleInfo {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	var rs GradleInfo
	gradleVer := regexp.MustCompile("Gradle\\s*([0-9A-Za-z_.-]+)")
	for _, line := range lines {
		if m := gradleVer.FindStringSubmatch(line); len(m) > 0 {
			rs.Version = m[1]
		}
		if strings.HasPrefix(s, "Revision: ") {
			rs.Revision = strings.TrimSpace(strings.TrimPrefix(line, "Revision: "))
		}
	}
	return rs
}

func execRawGradleInfo(baseDir string) (*GradleInfo, error) {
	c := exec.Command("gradle", "--version")
	c.Dir = baseDir
	data, e := c.Output()
	if e != nil {
		s := strings.TrimSpace(string(data))
		if len(s) > 64 {
			s = s[:64]
		}
		return nil, errors.Wrap(e, "Get version failed: "+s)
	}
	rs := parseGradleVersion(string(data))
	rs.CallCmd = func(args ...string) *exec.Cmd {
		return exec.Command("gradle", args...)
	}
	rs.Executable = "gradle"
	return &rs, nil
}

func execWrappedGradleInfo(baseDir string) (*GradleInfo, error) {
	var c *exec.Cmd
	var wrapperPath string
	if runtime.GOOS == "windows" {
		wrapperPath = filepath.Join(baseDir, "gradlew.bat")
		c = exec.Command(wrapperPath, "--version")
	} else {
		wrapperPath = filepath.Join(baseDir, "gradlew")
		d, e := exec.Command("chmod", "0755", wrapperPath).Output()
		if e != nil {
			logger.Warn.Println("Chmod wrapper 0755 failed.", e.Error(), string(d), wrapperPath)
		}
		c = exec.Command("sh", "-c", wrapperPath, "--version")
	}
	c.Dir = baseDir
	data, e := c.Output()
	if e != nil {
		// truncate output string if too long
		s := strings.TrimSpace(string(data))
		if len(s) > 64 {
			s = s[:64]
		}
		return nil, errors.Wrap(e, "Get version failed: "+s)
	}
	rs := parseGradleVersion(string(data))
	rs.Executable = wrapperPath
	if runtime.GOOS == "windows" {
		rs.CallCmd = func(args ...string) *exec.Cmd {
			return exec.Command(wrapperPath, args...)
		}
	} else {
		rs.CallCmd = func(args ...string) *exec.Cmd {
			aa := []string{wrapperPath}
			aa = append(aa, args...)
			return exec.Command("sh", "-c", strings.Join(aa, " "))
		}
	}
	return &rs, nil
}
