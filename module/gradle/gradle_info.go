package gradle

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/pkg/errors"
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
	CallCmd    func(ctx context.Context, args ...string) *exec.Cmd
}

func (g *GradleInfo) String() string {
	return fmt.Sprintf("Gradle[%s]: %s , revision: %s", g.Version, g.Executable, g.Revision)
}

func evalGradleInfo(ctx context.Context, dir string) (info *GradleInfo, e error) {
	gradlewDir := dir
	info, e = execWrappedGradleInfo(ctx, gradlewDir)
	if e == nil {
		return // gradle wrapper 找到了，就他了
	} else {
		logger.Debug.Println("check gradle wrapper failed.", e.Error())
	}
	info, e = execRawGradleInfo(ctx, dir)
	if e != nil {
		logger.Debug.Println("check raw gradle failed.", e.Error())
	}
	return
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

func execRawGradleInfo(ctx context.Context, baseDir string) (*GradleInfo, error) {
	c := exec.CommandContext(ctx, "gradle", "--version")
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
	rs.CallCmd = func(ctx context.Context, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "gradle", args...)
	}
	rs.Executable = "gradle"
	return &rs, nil
}

func execWrappedGradleInfo(ctx context.Context, baseDir string) (*GradleInfo, error) {
	var c *exec.Cmd
	var wrapperPath string
	if runtime.GOOS == "windows" {
		wrapperPath = filepath.Join(baseDir, "gradlew.bat")
		c = exec.CommandContext(ctx, wrapperPath, "--version", "--quiet")
	} else {
		wrapperPath = filepath.Join(baseDir, "gradlew")
		d, e := exec.Command("chmod", "0755", wrapperPath).Output()
		if e != nil {
			logger.Warn.Println("Chmod wrapper 0755 failed.", e.Error(), string(d), wrapperPath)
		}
		c = exec.CommandContext(ctx, wrapperPath, "--version", "--quiet")
	}
	logger.Debug.Println("Query version:", c.String())
	c.Dir = baseDir
	data, e := c.Output()
	if e != nil {
		// truncate output string if too long
		s := strings.TrimSpace(string(data))
		if len(s) > 1024 {
			s = s[:1024]
		}
		return nil, errors.Wrap(e, "Get version failed: "+s)
	}
	rs := parseGradleVersion(string(data))
	rs.Executable = wrapperPath
	if runtime.GOOS == "windows" {
		rs.CallCmd = func(ctx context.Context, args ...string) *exec.Cmd {
			c := exec.CommandContext(ctx, wrapperPath, args...)
			logger.Debug.Println("Execute:", c.String())
			return c
		}
	} else {
		rs.CallCmd = func(ctx context.Context, args ...string) *exec.Cmd {
			aa := []string{wrapperPath}
			aa = append(aa, args...)
			c := exec.CommandContext(ctx, "sh", "-c", strings.Join(aa, " "))
			logger.Debug.Println("Execute:", c.String())
			return c
		}
	}
	return &rs, nil
}
