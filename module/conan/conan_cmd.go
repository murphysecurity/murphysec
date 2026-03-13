package conan

import (
	"bytes"
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/logpipe"
	"github.com/murphysecurity/murphysec/infra/suffixbuf"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type CmdInfo struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

func (c CmdInfo) String() string {
	return fmt.Sprintf("[%s]%s", c.Version, c.Path)
}

var _conanCmdInfo any

func getConanInfo(ctx context.Context) (*CmdInfo, error) {
	if info, ok := _conanCmdInfo.(*CmdInfo); ok {
		return info, nil
	}
	if e, ok := _conanCmdInfo.(error); ok {
		return nil, e
	}
	p, e := LocateConan(ctx)
	if e != nil {
		_conanCmdInfo = e
		return nil, e
	}
	v, e := GetConanVersion(ctx, p)
	if e != nil {
		_conanCmdInfo = e
		return nil, e
	}
	r := &CmdInfo{
		Path:    p,
		Version: v,
	}
	_conanCmdInfo = r
	return r, nil
}

type ConanJsonKind string

const (
	ConanJsonKindGraph ConanJsonKind = "graph"
	ConanJsonKindInfo  ConanJsonKind = "info"
)

func logConanRemoteConfigPaths(logger *zap.Logger, major int) {
	const maxLogBytes = 64 * 1024
	home := os.Getenv("HOME")
	if home == "" {
		return
	}
	candidates := []string{
		filepath.Join(home, ".conan", "remotes.json"),
		filepath.Join(home, ".conan2", "remotes.json"),
	}
	for _, p := range candidates {
		_, err := os.Stat(p)
		exists := err == nil
		logger.Sugar().Infof("Conan remote config path: %s (exists=%t)", p, exists)
		if !exists {
			continue
		}
		data, readErr := os.ReadFile(p)
		if readErr != nil {
			logger.Sugar().Warnf("Conan remote config read failed: %s, err=%v", p, readErr)
			continue
		}
		content := string(data)
		if len(content) > maxLogBytes {
			content = content[:maxLogBytes] + "\n...(truncated)"
		}
		logger.Sugar().Infof("Conan remote config content (%s):\n%s", p, content)
	}
	if major >= 2 {
		logger.Info("Conan remote config in use: ~/.conan2/remotes.json (expected)")
	} else {
		logger.Info("Conan remote config in use: ~/.conan/remotes.json (expected)")
	}
}

func ExecuteConanInfoCmd(ctx context.Context, cmdInfo *CmdInfo, dir string) (string, ConanJsonKind, error) {
	logger := logctx.Use(ctx)
	lp := logpipe.New(logger, "conan")
	defer lp.Close()
	if cmdInfo == nil {
		return "", "", fmt.Errorf("conan cmd info is nil")
	}
	jsonP := getConanInfoJsonPath()
	major := ConanMajorVersion(cmdInfo.Version)
	logger.Sugar().Infof("Conan detected: path=%s version=%s major=%d", cmdInfo.Path, cmdInfo.Version, major)
	logConanRemoteConfigPaths(logger, major)
	logger.Sugar().Debugf("temp file: %s", jsonP)
	if major >= 2 {
		logger.Info("Conan mode selected: graph")
		if e := executeConanGraphInfoCmd(ctx, cmdInfo.Path, dir, jsonP); e != nil {
			return "", "", e
		}
		logger.Info("Conan command completed with graph mode")
		return jsonP, ConanJsonKindGraph, nil
	}
	logger.Info("Conan mode selected: info")
	if e := executeConanInfoCmd(ctx, cmdInfo.Path, dir, jsonP); e != nil {
		return "", "", e
	}
	logger.Info("Conan command completed with info mode")
	return jsonP, ConanJsonKindInfo, nil
}

func executeConanInfoCmd(ctx context.Context, conanPath string, dir string, jsonP string) error {
	logger := logctx.Use(ctx)
	c := exec.Command(conanPath, "info", ".", "-j", jsonP)
	logger.Sugar().Infof("Command: %s", c.String())
	c.Env = getEnvForConan()
	c.Dir = dir
	sb := suffixbuf.NewSize(1024)
	logPipe := logpipe.New(logger, "conan")
	defer logPipe.Close()
	c.Stdout = io.MultiWriter(sb, logPipe)
	c.Stderr = io.MultiWriter(sb, logPipe)
	if e := c.Run(); e != nil {
		logger.Warn("Conan command exit with error", zap.Error(e))
		return conanError(sb.Bytes())
	}
	return nil
}

func executeConanGraphInfoCmd(ctx context.Context, conanPath string, dir string, jsonP string) error {
	logger := logctx.Use(ctx)
	c := exec.Command(conanPath, "graph", "info", ".", "--format=json")
	logger.Sugar().Infof("Command: %s", c.String())
	c.Env = getEnvForConan()
	c.Dir = dir
	sb := suffixbuf.NewSize(1024)
	var out bytes.Buffer
	logPipe := logpipe.New(logger, "conan")
	defer logPipe.Close()
	c.Stdout = io.MultiWriter(&out, logPipe, sb)
	c.Stderr = io.MultiWriter(logPipe, sb)
	if e := c.Run(); e != nil {
		logger.Warn("Conan graph command exit with error", zap.Error(e))
		return conanError(sb.Bytes())
	}
	if e := os.WriteFile(jsonP, out.Bytes(), 0o644); e != nil {
		return fmt.Errorf("write conan graph json failed: %w", e)
	}
	return nil
}

func getConanInfoJsonPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprint("conan-info-", rand.Uint64(), ".json"))
}

func LocateConan(ctx context.Context) (string, error) {
	s, e := exec.LookPath("conan")
	if e != nil {
		return "", errors.WithCause(ErrConanNotFound, e)
	}
	return s, nil
}

func GetConanVersion(ctx context.Context, conanPath string) (string, error) {
	c := exec.CommandContext(ctx, conanPath, "-v")
	c.Env = getEnvForConan()
	if data, e := c.Output(); e != nil {
		return "", errors.WithCause(ErrGetConanVersionFail, e)
	} else {
		return strings.TrimSpace(strings.TrimPrefix(string(data), "Conan version")), nil
	}
}

func ConanMajorVersion(version string) int {
	version = strings.TrimSpace(version)
	if version == "" {
		return 0
	}
	parts := strings.SplitN(version, ".", 2)
	v, _ := strconv.Atoi(parts[0])
	return v
}

func getEnvForConan() []string {
	osEnv := os.Environ()
	var rs = make([]string, 0, len(osEnv)+3)
	rs = append(rs, osEnv...)
	return append(rs, "CONAN_NON_INTERACTIVE=1", "NO_COLOR=1", "CLICOLOR=0")
}
