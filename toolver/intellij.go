package toolver

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/murphysecurity/murphysec/utils"
)

var intellijConfig Config
var intellijConfigOnce sync.Once

func getIntellijConfig() *Config {
	intellijConfigOnce.Do(func() {
		var config = &intellijConfig
		config.Maven.MavenCommand = locateMavenIdeaMavenHome()
		if s := strings.TrimSpace(os.Getenv("IDEA_MAVEN_JRE")); s != "" && utils.IsDir(s) {
			config.Maven.JavaHome = s
		}
		if s := strings.TrimSpace(os.Getenv("IDEA_MAVEN_CONF")); s != "" && utils.IsFile(s) && strings.ToLower(filepath.Ext(s)) == ".xml" {
			config.Maven.MavenSettingPath = s
		}
	})
	return &intellijConfig
}

func locateMavenIdeaMavenHome() string {
	var s = strings.TrimSpace(os.Getenv("IDEA_MAVEN_HOME"))
	if s == "" {
		return ""
	}
	if !filepath.IsAbs(s) {
		abs, e := filepath.Abs(s)
		if e == nil {
			s = abs
		}
	}
	if _s, e := filepath.EvalSymlinks(s); e == nil && _s != "" {
		s = _s
	}
	if runtime.GOOS == "windows" {
		return locateMavenIdeaMavenHomeWindows(s)
	}
	return locateMavenIdeaMavenHomeUnix(s)
}

func locateMavenIdeaMavenHomeUnix(s string) string {
	if utils.IsFile(s) {
		return s
	}
	suffixes := []string{"mvn", "bin/mvn"}
	for _, it := range suffixes {
		target := filepath.Join(s, it)
		if utils.IsFile(target) {
			return target
		}
	}
	return ""
}

func locateMavenIdeaMavenHomeWindows(s string) string {
	if utils.IsFile(s) {
		ext := strings.ToLower(filepath.Ext(s))
		if ext == ".exe" || ext == ".bat" || ext == ".cmd" {
			return s
		}
		return ""
	}
	var suffixes = []string{
		"mvn.cmd",
		"mvn.bat",
		"mvn.exe",
		"bin\\mvn.cmd",
		"bin\\mvn.bat",
		"bin\\mvn.exe",
	}
	for _, suffix := range suffixes {
		target := filepath.Join(s, suffix)
		if utils.IsFile(target) {
			return target
		}
	}
	return ""
}
