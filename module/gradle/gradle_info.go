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

func evalGradleInfo(dir string) (info *GradleInfo, e error) {
	/**
	1.
	由于 ../../inspector/managed_inspect.go:35 使用了filepath.WalkDir
	这里必然是/xx/current，所以 filepath.Split 方法必然返回dir和name

	2.
	根据gradle wrapper的介绍，执行gradle wrapper命令后，目录结构如下：
	.
	├── a-subproject
	│   └── build.gradle
	├── settings.gradle
	├── gradle
	│   └── wrapper
	│       ├── gradle-wrapper.jar
	│       └── gradle-wrapper.properties
	├── gradlew
	└── gradlew.bat
	所以，需要向上一级，才能正确访问gradlew/gradlew.bat
	https://docs.gradle.org/current/userguide/gradle_wrapper.html

	==== 此类目录结构并非强制，但处于兼容考虑应当支持~~~
	*/
	gradlewDir := dir
	for backTrackCount := 0; backTrackCount < 2 && gradlewDir != ""; backTrackCount++ {
		info, e = execWrappedGradleInfo(gradlewDir)
		if e == nil {
			return // gradle wrapper 找到了，就他了
		} else {
			logger.Debug.Println("check gradle wrapper failed.", e.Error())
			gradlewDir = filepath.Dir(gradlewDir)
		}
	}
	info, e = execRawGradleInfo(dir)
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
		c = exec.Command(wrapperPath, "--version")
	}
	logger.Debug.Println("Query version:", c.String())
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
			c := exec.Command(wrapperPath, args...)
			logger.Debug.Println("Execute:", c.String())
			return c
		}
	} else {
		rs.CallCmd = func(args ...string) *exec.Cmd {
			aa := []string{wrapperPath}
			aa = append(aa, args...)
			c := exec.Command("sh", "-c", strings.Join(aa, " "))
			logger.Debug.Println("Execute:", c.String())
			return c
		}
	}
	return &rs, nil
}
