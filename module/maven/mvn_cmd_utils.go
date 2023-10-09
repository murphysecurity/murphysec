package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const M2SettingsFilePathCtxKey = "MavenSettingsFilePathCtxKey"

type MvnCommandInfo struct {
	Path             string `json:"path"`
	MvnVersion       string `json:"mvn_version"`
	UserSettingsPath string `json:"user_settings_path"`
	JavaHome         string `json:"java_home"`
}

func (m MvnCommandInfo) String() string {
	return fmt.Sprintf("MavenCommand: %s, JavaHome: %s, MavenVersion: %s, UserSettings: %s", m.Path, m.JavaHome, m.MvnVersion, m.UserSettingsPath)
}

func (m MvnCommandInfo) Command(ctx context.Context, args ...string) *exec.Cmd {
	if ctx == nil {
		ctx = context.TODO()
	}
	var _args = make([]string, 0, len(args)+5)
	if m.UserSettingsPath != "" {
		_args = append(_args, "--settings", m.UserSettingsPath)
	}
	_args = append(_args, "--batch-mode")
	_args = append(_args, args...)
	cmd := exec.CommandContext(ctx, m.Path, _args...)
	if m.JavaHome != "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "JAVA_HOME="+m.JavaHome)
	}
	return cmd
}

var cachedMvnCommandResult *_MvnCommandResult

type _MvnCommandResult struct {
	rs *MvnCommandInfo
	e  error
}

func CheckMvnCommand(ctx context.Context) (info *MvnCommandInfo, err error) {
	var logger = logctx.Use(ctx)
	if cachedMvnCommandResult != nil {
		if cachedMvnCommandResult.e != nil {
			logger.Warn("Cached maven error", zap.Error(cachedMvnCommandResult.e))
		}
		if cachedMvnCommandResult.rs != nil {
			logger.Debug("Use cached maven command info", zap.String("info", cachedMvnCommandResult.rs.String()))
		}
		return cachedMvnCommandResult.rs, cachedMvnCommandResult.e
	}
	defer func() {
		cachedMvnCommandResult = &_MvnCommandResult{
			rs: info,
			e:  err,
		}
	}()
	if env.DisableMvnCommand {
		return nil, ErrMvnDisabled.Detailed("environment variable NO_MVN set")
	}
	info = &MvnCommandInfo{}
	info.Path = env.IdeaMavenHome
	if info.Path == "" {
		info.Path = getMvnCommandOs()
	}
	if info.Path == "" {
		return nil, ErrMvnNotFound
	}
	info.JavaHome = env.IdeaMavenJre
	info.UserSettingsPath = env.IdeaMavenConf
	if r, ok := ctx.Value(M2SettingsFilePathCtxKey).(string); ok {
		info.UserSettingsPath = r
	}
	// check version
	ver, e := checkMvnVersion(ctx, info.Path, info.JavaHome)
	if e != nil {
		return nil, e
	}
	info.MvnVersion = ver
	return
}

func locateMvnCmdPath() string {
	return getMvnCommandOs()
}

func executeMvnVersion(ctx context.Context, mvnPath string, javaHome string) (string, error) {
	var logger = logctx.Use(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*8)
	defer cancel()
	cmd := exec.CommandContext(ctx, mvnPath, "--version", "--batch-mode")
	cmd.Env = os.Environ()
	if javaHome != "" {
		cmd.Env = append(cmd.Env, "JAVA_HOME="+javaHome)
	}
	logger.Info("Check maven version", zap.String("cmd", cmd.String()))
	output, err := cmd.Output()
	if err != nil {
		logger.Error("Check maven version failed", zap.Error(err))
		return "", ErrCheckMvnVersion.Wrap(err)
	}
	return string(output), nil
}

func checkMvnVersion(ctx context.Context, mvnPath string, javaHome string) (string, error) {
	var logger = logctx.Use(ctx)
	output, err := executeMvnVersion(ctx, mvnPath, javaHome)
	if err != nil {
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			logger.Warn("System is linux or darwin, try to grant executable permission")
			_ = os.Chmod(mvnPath, 0755)
			logger.Warn("Retry...")
			output, err = executeMvnVersion(ctx, mvnPath, javaHome)
		}
		if err != nil {
			return "", err
		}
	}
	ver := parseMvnVersion(output)
	if ver == "" {
		return "", ErrCheckMvnVersion
	}
	return ver, nil
}

func parseMvnVersion(input string) string {
	versionPattern := regexp.MustCompile(`Apache Maven (\d+(?:\.[\dA-Za-z_-]+)+)`)
	lines := strings.Split(input, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	for _, it := range lines {
		line := strings.TrimSpace(it)
		if m := versionPattern.FindStringSubmatch(line); m != nil {
			return m[1]
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
