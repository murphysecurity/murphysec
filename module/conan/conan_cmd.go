package conan

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/logpipe"
	"github.com/murphysecurity/murphysec/infra/suffixbuf"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CmdInfo struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

func (c CmdInfo) String() string {
	return fmt.Sprintf("[%s]%s", c.Version, c.Path)
}

var _conanCmdInfo any

const conanVerboseArg = "debug"

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
		truncated := false
		if len(data) > maxLogBytes {
			data = data[:maxLogBytes]
			truncated = true
		}
		contentBase64 := base64.StdEncoding.EncodeToString(data)
		if truncated {
			contentBase64 += "\n...(truncated raw bytes)"
		}
		logger.Sugar().Infof("Conan remote config content base64 (%s):\n%s", p, contentBase64)
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
	remoteCreds, credErr := getConanRemoteCredentialsFromConfig(major)
	if credErr != nil {
		logger.Warn("Conan remote credential precheck failed when reading config", zap.Error(credErr))
	}
	logger.Sugar().Infof("Conan detected: path=%s version=%s major=%d", cmdInfo.Path, cmdInfo.Version, major)
	logger.Sugar().Infof("Conan verbose mode: -v %s", conanVerboseArg)
	logConanRemoteConfigPaths(logger, major)
	loginWarnings := ensureConanRemoteLogin(ctx, cmdInfo.Path, major, remoteCreds)
	logger.Sugar().Debugf("temp file: %s", jsonP)
	if major >= 2 {
		if e := ensureConan2DefaultProfile(ctx, cmdInfo.Path, major, remoteCreds); e != nil {
			return "", "", e
		}
		logger.Info("Conan mode selected: graph")
		if e := executeConanGraphInfoCmd(ctx, cmdInfo.Path, dir, jsonP, major, remoteCreds); e != nil {
			if len(loginWarnings) > 0 {
				return "", "", fmt.Errorf("%w; login warnings: %s", e, strings.Join(loginWarnings, " | "))
			}
			return "", "", e
		}
		logger.Info("Conan command completed with graph mode")
		return jsonP, ConanJsonKindGraph, nil
	}
	logger.Info("Conan mode selected: info")
	if e := executeConanInfoCmd(ctx, cmdInfo.Path, dir, jsonP, major, remoteCreds); e != nil {
		if len(loginWarnings) > 0 {
			return "", "", fmt.Errorf("%w; login warnings: %s", e, strings.Join(loginWarnings, " | "))
		}
		return "", "", e
	}
	logger.Info("Conan command completed with info mode")
	return jsonP, ConanJsonKindInfo, nil
}

func ensureConan2DefaultProfile(ctx context.Context, conanPath string, major int, remoteCreds []conanRemoteCredential) error {
	logger := logctx.Use(ctx)
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		logger.Warn("Conan profile precheck skipped: cannot determine home directory", zap.Error(err))
		return nil
	}
	profilePath := filepath.Join(home, ".conan2", "profiles", "default")
	if _, statErr := os.Stat(profilePath); statErr == nil {
		logger.Sugar().Infof("Conan default profile exists: %s", profilePath)
		return nil
	} else if !os.IsNotExist(statErr) {
		logger.Sugar().Warnf("Conan default profile stat failed: %s, err=%v", profilePath, statErr)
	}

	logger.Sugar().Infof("Conan default profile missing: %s, running detect", profilePath)
	args := conanArgs(major, "profile", "detect", "--force")
	c := exec.CommandContext(ctx, conanPath, args...)
	logger.Sugar().Infof("Command: %s", c.String())
	c.Env = getEnvForConan(major, remoteCreds)
	start := time.Now()
	sb := suffixbuf.NewSize(1024)
	logPipe := logpipe.New(logger, "conan")
	defer logPipe.Close()
	c.Stdout = io.MultiWriter(sb, logPipe)
	c.Stderr = io.MultiWriter(sb, logPipe)
	if e := c.Run(); e != nil {
		logger.Warn("Conan profile detect command exit with error", zap.Error(e))
		return conanError(sb.Bytes())
	}
	logger.Sugar().Infof("Conan profile detect completed in %s", time.Since(start))

	if _, statErr := os.Stat(profilePath); statErr != nil {
		return fmt.Errorf("conan profile detect completed but default profile still missing: %s, err=%w", profilePath, statErr)
	}
	logger.Sugar().Infof("Conan default profile created: %s", profilePath)
	return nil
}

func executeConanInfoCmd(ctx context.Context, conanPath string, dir string, jsonP string, major int, remoteCreds []conanRemoteCredential) error {
	logger := logctx.Use(ctx)
	args := conanArgs(major, "info", ".", "-j", jsonP)
	c := exec.Command(conanPath, args...)
	logger.Sugar().Infof("Command: %s", c.String())
	c.Env = getEnvForConan(major, remoteCreds)
	c.Dir = dir
	start := time.Now()
	sb := suffixbuf.NewSize(1024)
	logPipe := logpipe.New(logger, "conan")
	defer logPipe.Close()
	c.Stdout = io.MultiWriter(sb, logPipe)
	c.Stderr = io.MultiWriter(sb, logPipe)
	if e := c.Run(); e != nil {
		logger.Warn("Conan command exit with error", zap.Error(e))
		return conanError(sb.Bytes())
	}
	logger.Sugar().Infof("Conan info command completed in %s", time.Since(start))
	return nil
}

func executeConanGraphInfoCmd(ctx context.Context, conanPath string, dir string, jsonP string, major int, remoteCreds []conanRemoteCredential) error {
	logger := logctx.Use(ctx)
	args := conanArgs(major, "graph", "info", ".", "--format=json")
	c := exec.Command(conanPath, args...)
	logger.Sugar().Infof("Command: %s", c.String())
	c.Env = getEnvForConan(major, remoteCreds)
	c.Dir = dir
	start := time.Now()
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
	logger.Sugar().Infof("Conan graph info command completed in %s", time.Since(start))
	if e := os.WriteFile(jsonP, out.Bytes(), 0o644); e != nil {
		return fmt.Errorf("write conan graph json failed: %w", e)
	}
	return nil
}

type conanRemoteCredential struct {
	Name     string
	Username string
	Password string
}

func ensureConanRemoteLogin(ctx context.Context, conanPath string, major int, creds []conanRemoteCredential) []string {
	warnings := make([]string, 0)
	if major < 2 {
		return warnings
	}
	logger := logctx.Use(ctx)
	if len(creds) == 0 {
		logger.Info("Conan remote login skipped: no credentials found in remote URLs")
		return warnings
	}
	for _, cred := range creds {
		args := conanArgs(major, "remote", "auth", cred.Name, "--force")
		c := exec.CommandContext(ctx, conanPath, args...)
		c.Env = getEnvForConan(major, []conanRemoteCredential{cred})
		sb := suffixbuf.NewSize(1024)
		logPipe := logpipe.New(logger, "conan")
		logger.Sugar().Infof("Command: %s remote auth %s --force (credentials from env)", conanPath, cred.Name)
		c.Stdout = io.MultiWriter(sb, logPipe)
		c.Stderr = io.MultiWriter(sb, logPipe)
		if runErr := c.Run(); runErr != nil {
			logPipe.Close()
			msg := fmt.Sprintf("conan remote auth failed for %s(user=%s): %v, details: %s", cred.Name, cred.Username, runErr, strings.TrimSpace(string(sb.Bytes())))
			logger.Warn(msg)
			warnings = append(warnings, msg)
			continue
		}
		logPipe.Close()
		logger.Sugar().Infof("Conan remote auth completed: remote=%s user=%s", cred.Name, cred.Username)
	}
	return warnings
}

type conanRemotesFile struct {
	Remotes []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"remotes"`
}

func getConanRemoteCredentialsFromConfig(major int) ([]conanRemoteCredential, error) {
	path, err := conanRemotesConfigPath(major)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg conanRemotesFile
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse remotes config failed: %w", err)
	}
	rs := make([]conanRemoteCredential, 0)
	seen := make(map[string]struct{})
	for _, remote := range cfg.Remotes {
		u, parseErr := url.Parse(strings.TrimSpace(remote.URL))
		if parseErr != nil || u.User == nil {
			continue
		}
		username := u.User.Username()
		password, ok := u.User.Password()
		if username == "" || !ok || password == "" {
			continue
		}
		key := remote.Name + "|" + username
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		rs = append(rs, conanRemoteCredential{
			Name:     remote.Name,
			Username: username,
			Password: password,
		})
	}
	return rs, nil
}

func conanRemotesConfigPath(major int) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", fmt.Errorf("cannot determine user home: %w", err)
	}
	if major >= 2 {
		return filepath.Join(home, ".conan2", "remotes.json"), nil
	}
	return filepath.Join(home, ".conan", "remotes.json"), nil
}

func conanArgs(major int, args ...string) []string {
	if major >= 2 {
		args = append(args, "-cc", "core:non_interactive=True")
	}
	return append(args, "-v", conanVerboseArg)
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
	c.Env = getBaseEnvForConan()
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

func getBaseEnvForConan() []string {
	osEnv := os.Environ()
	var rs = make([]string, 0, len(osEnv)+3)
	rs = append(rs, osEnv...)
	return append(rs, "CONAN_NON_INTERACTIVE=1", "NO_COLOR=1", "CLICOLOR=0")
}

func getEnvForConan(major int, remoteCreds []conanRemoteCredential) []string {
	rs := getBaseEnvForConan()
	if major < 1 || len(remoteCreds) == 0 {
		return rs
	}
	for _, cred := range remoteCreds {
		suffix := conanRemoteEnvVarSuffix(cred.Name)
		if suffix == "" || cred.Username == "" || cred.Password == "" {
			continue
		}
		rs = append(rs,
			fmt.Sprintf("CONAN_LOGIN_USERNAME_%s=%s", suffix, cred.Username),
			fmt.Sprintf("CONAN_PASSWORD_%s=%s", suffix, cred.Password),
		)
	}
	return rs
}

func conanRemoteEnvVarSuffix(remoteName string) string {
	remoteName = strings.TrimSpace(remoteName)
	if remoteName == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range remoteName {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r - ('a' - 'A'))
		case (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}
