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

func ExecuteConanInfoCmd(ctx context.Context, conanPath string, dir string) (string, ConanJsonKind, error) {
	logger := logctx.Use(ctx)
	lp := logpipe.New(logger, "conan")
	defer lp.Close()
	jsonP := getConanInfoJsonPath()
	logger.Sugar().Debugf("temp file: %s", jsonP)
	// Prefer Conan2 graph mode for better compatibility with modern recipes.
	if e := executeConanGraphInfoCmd(ctx, conanPath, dir, jsonP); e == nil {
		logger.Info("Conan command completed with graph mode")
		return jsonP, ConanJsonKindGraph, nil
	}
	logger.Sugar().Warn("Conan graph mode failed, falling back to info mode")

	if e := executeConanInfoCmd(ctx, conanPath, dir, jsonP); e != nil {
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
