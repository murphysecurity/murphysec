package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/utils"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type MvnCommandInfo struct {
	Path       string `json:"path"`
	MvnVersion string `json:"mvn_version"`
}

func (m MvnCommandInfo) String() string {
	return fmt.Sprintf("[%s]%s", m.MvnVersion, m.Path)
}

var cachedMvnCommandResult *_MvnCommandResult

type _MvnCommandResult struct {
	rs *MvnCommandInfo
	e  error
}

func CheckMvnCommand() (info *MvnCommandInfo, err error) {
	if cachedMvnCommandResult != nil {
		return cachedMvnCommandResult.rs, cachedMvnCommandResult.e
	}
	defer func() {
		cachedMvnCommandResult = &_MvnCommandResult{
			rs: info,
			e:  err,
		}
	}()
	if env.DisableMvnCommand {
		return nil, ErrMvnDisabled
	}
	mvnPath := getMvnCommandPath()
	if mvnPath == "" {
		return nil, ErrMvnNotFound
	}
	// check version
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	output, err := exec.CommandContext(ctx, mvnPath, "--version").Output()
	if err != nil {
		return nil, ErrCheckMvnVersion.Wrap(err)
	}
	versionPattern := regexp.MustCompile("Apache Maven (\\d+(?:\\.[\\dA-Za-z_-]+)+)")
	lines := strings.Split(string(output), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	for _, it := range lines {
		line := strings.TrimSpace(it)
		if m := versionPattern.FindStringSubmatch(line); m != nil {
			return &MvnCommandInfo{
				Path:       mvnPath,
				MvnVersion: m[1],
			}, nil
		}
	}
	return nil, ErrCheckMvnVersion
}

func getMvnCommandPath() string {
	if f := getMvnCommandOs(); f != "" {
		return f
	}
	return getMvnCommandIntellijIDEA()
}

func getMvnCommandIntellijIDEA() string {
	if env.IdeaInstallPath == "" {
		return ""
	}
	var name []string
	name = append(name, "mvn")
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		name = append(name, "mvn.cmd", "mvn.bat")
	} else {
		name = append(name, "mvn.sh")
	}
	var refPath []string
	refPath = append(refPath, "plugins/maven/lib/maven3/bin")

	for _, ref := range refPath {
		for _, n := range name {
			p := filepath.Join(env.IdeaInstallPath, ref, n)
			if !filepath.IsAbs(p) {
				var e error
				p, e = filepath.Abs(p)
				if e != nil {
					continue
				}
			}
			if utils.IsFile(filepath.Join(ref, n)) {
				return p
			}
		}
	}
	return ""
}

func getMvnCommandOs() string {
	p, e := exec.LookPath("mvn")
	if e != nil {
		return ""
	}
	if filepath.IsAbs(p) {
		return p
	}
	if p, e := filepath.Abs(p); e == nil {
		return p
	}
	return ""
}
