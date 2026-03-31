package nuget

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/murphysecurity/murphysec/utils"

	"sync"

	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/logsanitize"
	"github.com/murphysecurity/murphysec/model"
	"go.uber.org/zap"
)

var _ErrDotnetNotFound = errors.New("dotnet not found")

const (
	nugetBuildMaxTimeout    = 6 * time.Hour
	nugetCommandIdleTimeout = 30 * time.Second
	nugetPreflightTimeout   = 8 * time.Second
	nugetRestoreBinlogDir   = ".murphysec"
	nugetTmpBinlogDir       = "murphysec-nuget-binlogs"
)

type nugetConfigXML struct {
	PackageSources nugetPackageSourcesXML `xml:"packageSources"`
}

type nugetPackageSourcesXML struct {
	Clear *struct{}           `xml:"clear"`
	Add   []nugetSourceAddXML `xml:"add"`
}

type nugetSourceAddXML struct {
	Value string `xml:"value,attr"`
}

type nugetSourceProbeResult struct {
	Source     string
	Reachable  bool
	StatusCode int
	Err        error
	Duration   time.Duration
}

func tailText(s string, max int) string {
	s = strings.TrimSpace(s)
	if s == "" || max <= 0 || len(s) <= max {
		return s
	}
	return "...(truncated)\n" + s[len(s)-max:]
}

func appendIfExists(paths []string, seen map[string]struct{}, p string) []string {
	if p == "" {
		return paths
	}
	if _, ok := seen[p]; ok {
		return paths
	}
	seen[p] = struct{}{}
	return append(paths, p)
}

func detectNugetConfigPaths(solutionPath string) []string {
	seen := map[string]struct{}{}
	var paths []string
	dir := filepath.Dir(solutionPath)
	for {
		paths = appendIfExists(paths, seen, filepath.Join(dir, "NuGet.Config"))
		paths = appendIfExists(paths, seen, filepath.Join(dir, "nuget.config"))
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	home := os.Getenv("HOME")
	if home != "" {
		paths = appendIfExists(paths, seen, filepath.Join(home, ".nuget", "NuGet", "NuGet.Config"))
		paths = appendIfExists(paths, seen, filepath.Join(home, ".config", "NuGet", "NuGet.Config"))
	}
	paths = appendIfExists(paths, seen, filepath.Join(string(os.PathSeparator), "etc", "nuget", "NuGet.Config"))
	paths = appendIfExists(paths, seen, filepath.Join(string(os.PathSeparator), "usr", "local", "share", "NuGet", "Config", "NuGet.Config"))
	return paths
}

func logNugetRemoteConfigPaths(logger *zap.Logger, solutionPath string) {
	const maxLogBytes = 64 * 1024
	for _, p := range detectNugetConfigPaths(solutionPath) {
		_, err := os.Stat(p)
		exists := err == nil
		logger.Sugar().Infof("NuGet remote config path: %s (exists=%t)", p, exists)
		if !exists {
			continue
		}
		data, readErr := os.ReadFile(p)
		if readErr != nil {
			logger.Sugar().Warnf("NuGet remote config read failed: %s, err=%v", p, readErr)
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
		logger.Sugar().Infof("NuGet remote config content base64 (%s):\n%s", p, contentBase64)
	}
}

func logNugetToolVersion(ctx context.Context, logger *zap.Logger) {
	dotnetPath, err := exec.LookPath("dotnet")
	if err != nil {
		logger.Warn("NuGet tool detect failed: dotnet not found", zap.Error(err))
		return
	}
	cmd := exec.CommandContext(ctx, dotnetPath, "--version")
	data, err := cmd.Output()
	if err != nil {
		logger.Warn("NuGet tool detect failed: dotnet --version failed", zap.String("path", dotnetPath), zap.Error(err))
		return
	}
	logger.Sugar().Infof("NuGet detected: path=%s version=%s", dotnetPath, strings.TrimSpace(string(data)))
}

func collectNugetSources(solutionPath string) ([]string, error) {
	paths := detectNugetConfigPaths(solutionPath)
	existing := make([]string, 0, len(paths))
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			existing = append(existing, p)
		}
	}
	// Apply low-priority first, then high-priority overrides.
	for i, j := 0, len(existing)-1; i < j; i, j = i+1, j-1 {
		existing[i], existing[j] = existing[j], existing[i]
	}

	sources := make([]string, 0)
	seen := map[string]struct{}{}
	for _, p := range existing {
		items, hasClear, err := parseNugetSourcesFromConfig(p)
		if err != nil {
			return nil, fmt.Errorf("parse NuGet config failed (%s): %w", p, err)
		}
		if hasClear {
			sources = sources[:0]
			seen = map[string]struct{}{}
		}
		for _, s := range items {
			if _, ok := seen[s]; ok {
				continue
			}
			seen[s] = struct{}{}
			sources = append(sources, s)
		}
	}
	return sources, nil
}

func parseNugetSourcesFromConfig(path string) ([]string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, err
	}
	var cfg nugetConfigXML
	if err = xml.Unmarshal(data, &cfg); err != nil {
		return nil, false, err
	}
	out := make([]string, 0, len(cfg.PackageSources.Add))
	for _, item := range cfg.PackageSources.Add {
		v := strings.TrimSpace(item.Value)
		if v != "" {
			out = append(out, v)
		}
	}
	return out, cfg.PackageSources.Clear != nil, nil
}

func maskURLCredential(raw string) string {
	u, err := neturl.Parse(raw)
	if err != nil {
		return raw
	}
	if u.User != nil {
		user := u.User.Username()
		if user != "" {
			u.User = neturl.UserPassword(user, "******")
		}
	}
	return u.String()
}

func probeNugetSource(ctx context.Context, source string) nugetSourceProbeResult {
	start := time.Now()
	res := nugetSourceProbeResult{Source: source}
	reqCtx, cancel := context.WithTimeout(ctx, nugetPreflightTimeout)
	defer cancel()

	transport := &http.Transport{}
	if os.Getenv("TLS_ALLOW_INSECURE") == "1" {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   nugetPreflightTimeout,
	}
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, source, nil)
	if err != nil {
		res.Err = err
		res.Duration = time.Since(start)
		return res
	}
	resp, err := client.Do(req)
	if err != nil {
		res.Err = err
		res.Duration = time.Since(start)
		return res
	}
	defer resp.Body.Close()
	_, _ = io.CopyN(io.Discard, resp.Body, 1024)
	res.StatusCode = resp.StatusCode
	res.Duration = time.Since(start)
	res.Reachable = (resp.StatusCode >= 200 && resp.StatusCode < 400) || resp.StatusCode == 401 || resp.StatusCode == 403
	return res
}

func preflightNugetSources(ctx context.Context, logger *zap.Logger, solutionPath string) error {
	sources, err := collectNugetSources(solutionPath)
	if err != nil {
		logger.Sugar().Warnf("NuGet source preflight parse failed, skip fast-fail: %v", err)
		return nil
	}
	if len(sources) == 0 {
		logger.Warn("NuGet source preflight skipped: no package sources found")
		return nil
	}
	logger.Sugar().Infof("NuGet source preflight start: total=%d", len(sources))
	results := make([]nugetSourceProbeResult, 0, len(sources))
	reachable := 0
	for _, s := range sources {
		r := probeNugetSource(ctx, s)
		results = append(results, r)
		if r.Reachable {
			reachable++
			logger.Sugar().Infof("NuGet source preflight ok: source=%s status=%d cost=%s",
				maskURLCredential(s), r.StatusCode, r.Duration.Round(time.Millisecond))
			continue
		}
		if r.Err != nil {
			logger.Sugar().Warnf("NuGet source preflight failed: source=%s err=%v cost=%s",
				maskURLCredential(s), r.Err, r.Duration.Round(time.Millisecond))
		} else {
			logger.Sugar().Warnf("NuGet source preflight failed: source=%s status=%d cost=%s",
				maskURLCredential(s), r.StatusCode, r.Duration.Round(time.Millisecond))
		}
	}
	if reachable > 0 {
		logger.Sugar().Infof("NuGet source preflight completed: reachable=%d/%d", reachable, len(sources))
		return nil
	}
	parts := make([]string, 0, len(results))
	for _, r := range results {
		if r.Err != nil {
			parts = append(parts, fmt.Sprintf("%s err=%v", maskURLCredential(r.Source), r.Err))
		} else {
			parts = append(parts, fmt.Sprintf("%s status=%d", maskURLCredential(r.Source), r.StatusCode))
		}
	}
	return fmt.Errorf("NuGet source preflight failed: all sources unreachable (%d). details: %s", len(results), strings.Join(parts, "; "))
}

func multipleBuilds(ctx context.Context, task *model.InspectionTask) error {
	logger := logctx.Use(ctx)
	slnPaths, err := findCLNList(task.Dir())
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	preferredSlnPaths, fallbackSlnPaths := splitSlnPathsByDockerCompose(slnPaths)
	targetSlnPaths := preferredSlnPaths
	if len(targetSlnPaths) == 0 {
		targetSlnPaths = fallbackSlnPaths
	}
	if err = validateBuildTargets(targetSlnPaths); err != nil {
		return err
	}
	logger.Sugar().Debugf("findCLNList: %v", slnPaths)
	logger.Sugar().Debugf("findCLNList preferred: %v", preferredSlnPaths)
	logger.Sugar().Debugf("findCLNList fallback: %v", fallbackSlnPaths)
	logger.Sugar().Debugf("findCLNList target: %v", targetSlnPaths)
	numCPU := utils.Coerce(runtime.NumCPU(), 1, 4)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	ch := make(chan string, len(targetSlnPaths))
	for _, j := range targetSlnPaths {
		ch <- j
	}
	close(ch)
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range ch {
				if err := buildEntrance(ctx, task, j); err != nil {
					logger.Warn(j + "buildEntrance faild:" + err.Error())
					mu.Lock()
					errs = append(errs, fmt.Errorf("%s: %w", j, err))
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil

}
func buildEntrance(ctx context.Context, task *model.InspectionTask, solutionPath string) error {
	logger := logctx.Use(ctx)
	if !isDotnetBuildTarget(solutionPath) {
		return fmt.Errorf("unsupported nuget build target: %s", solutionPath)
	}
	if missingRefs, e := findMissingProjectReferences(solutionPath); e == nil && len(missingRefs) > 0 {
		logger.Sugar().Warnf("skip nuget build for %s because referenced projects are missing in current scan context: %v", solutionPath, missingRefs)
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, nugetBuildMaxTimeout)
	defer cancel()
	e := listNuget(ctx, task, solutionPath)
	if e != nil {
		if errors.Is(e, _ErrDotnetNotFound) {
			logger.Warn("Dotnet not found, skip DotnetList")
			return e
		} else {
			// log it and go on
			logger.Warn("Dotnet list failed"+solutionPath, zap.Error(e))
			return e
		}
	}

	// 处理超时的情况
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			logger.Warn("restore timed out")
			return errors.New("restore timed out")
		}
		return errors.New("restore timed out")
	default:
		return nil
	}

}
func readOutput(pipe io.ReadCloser, logger *zap.Logger, logPrefix string) string {
	return readOutputWithHook(pipe, logger, logPrefix, nil)
}

func readOutputWithHook(pipe io.ReadCloser, logger *zap.Logger, logPrefix string, onLine func()) string {
	var res strings.Builder
	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(nil, 1024*1024)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if onLine != nil {
			onLine()
		}
		if logger != nil {
			logger.Debug(logPrefix + logsanitize.ForLog(line))
		}
		res.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil && logger != nil {
		logger.Error("reading output failed", zap.Error(err))
	}
	return res.String()
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n")
	if strings.HasSuffix(s, "\n") {
		return n
	}
	return n + 1
}

func startCommandIdleWatchdog(
	ctx context.Context,
	logger *zap.Logger,
	cmd *exec.Cmd,
	commandName string,
	target string,
	idleTimeout time.Duration,
	idleTimeoutExceeded *atomic.Bool,
	lastOutputAt *atomic.Int64,
) func() {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				last := time.Unix(0, lastOutputAt.Load())
				if time.Since(last) < idleTimeout {
					continue
				}
				idleTimeoutExceeded.Store(true)
				logger.Sugar().Warnf("%s idle timeout reached: no output for %s, target=%s, killing process",
					commandName, idleTimeout, target)
				if cmd.Process != nil {
					_ = cmd.Process.Kill()
				}
				return
			}
		}
	}()
	return func() {
		close(done)
	}
}

func nugetRestoreBinlogName(solutionPath string) string {
	base := filepath.Base(solutionPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if name == "" {
		name = "restore"
	}
	return name + ".restore.binlog"
}

func nugetRestoreBinlogPath(solutionPath string) string {
	return filepath.Join(filepath.Dir(solutionPath), nugetRestoreBinlogDir, nugetRestoreBinlogName(solutionPath))
}

func nugetRestoreTmpBinlogPath(solutionPath string) string {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(filepath.Clean(solutionPath)))
	return filepath.Join(
		os.TempDir(),
		nugetTmpBinlogDir,
		fmt.Sprintf("%08x", hasher.Sum32()),
		nugetRestoreBinlogName(solutionPath),
	)
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return dst.Close()
}

func archiveNugetRestoreBinlog(logger *zap.Logger, solutionPath string) (string, error) {
	srcPath := nugetRestoreBinlogPath(solutionPath)
	if _, err := os.Stat(srcPath); err != nil {
		return "", err
	}
	dstPath := nugetRestoreTmpBinlogPath(solutionPath)
	if err := copyFile(srcPath, dstPath); err != nil {
		return "", err
	}
	if logger != nil {
		logger.Sugar().Infof("NuGet restore binary log copied to tmp: %s", dstPath)
	}
	return dstPath, nil
}

func dotnetRestoreArgs(solutionPath string) []string {
	binlogPath := nugetRestoreBinlogPath(solutionPath)
	args := []string{
		"restore",
		solutionPath,
		"-v",
		"detailed",
		fmt.Sprintf("/bl:%s;ProjectImports=Embed", binlogPath),
	}
	if runtime.GOOS == "linux" {
		args = append(args, "-p:EnableWindowsTargeting=true")
	}
	return args
}

func dotnetListPackageArgs(solutionPath string) []string {
	return []string{"list", solutionPath, "package", "--include-transitive", "--format", "json"}
}

// 通过先运行 dotnet restore 命令，确保项目中的所有 NuGet 包依赖项被正确恢复
func buildPackage(ctx context.Context, logger *zap.Logger, solutionPath string) (err error) {
	//dotnet restore
	binlogPath := nugetRestoreBinlogPath(solutionPath)
	if err = os.MkdirAll(filepath.Dir(binlogPath), 0o755); err != nil {
		err = fmt.Errorf("create nuget binlog dir failed: %w", err)
		logger.Error(err.Error())
		return
	}
	args := dotnetRestoreArgs(solutionPath)
	logger.Sugar().Infof("NuGet restore binary log enabled: %s", binlogPath)
	if runtime.GOOS == "linux" {
		logger.Info("dotnet restore adds EnableWindowsTargeting for Linux compatibility")
	}
	cmd := exec.CommandContext(ctx, "dotnet", args...)
	cmd.Dir = filepath.Dir(solutionPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("create stdout pipe failed: %w", err)
		logger.Error(err.Error())
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = fmt.Errorf("create stderr pipe failed: %w", err)
		logger.Error(err.Error())
		return
	}
	startAt := time.Now()
	lastOutputAt := atomic.Int64{}
	lastOutputAt.Store(startAt.UnixNano())
	idleTimeoutExceeded := atomic.Bool{}

	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("start command failed: %v", err)
		logger.Error(err.Error())
		return
	}
	stopWatchdog := startCommandIdleWatchdog(
		ctx, logger, cmd, "dotnet restore", solutionPath, nugetCommandIdleTimeout, &idleTimeoutExceeded, &lastOutputAt,
	)
	defer stopWatchdog()

	var errOutput strings.Builder
	var stderrOutput strings.Builder
	var stdoutOutput strings.Builder
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				last := time.Unix(0, lastOutputAt.Load())
				logger.Sugar().Infof("dotnet restore is still running: elapsed=%s last_output_ago=%s target=%s",
					time.Since(startAt).Round(time.Second),
					time.Since(last).Round(time.Second),
					solutionPath,
				)
			}
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	var stderrWg sync.WaitGroup
	stderrWg.Add(1)
	go func() {
		defer stderrWg.Done()
		scanner := bufio.NewScanner(stderr)
		scanner.Buffer(nil, 1024*4)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			lastOutputAt.Store(time.Now().UnixNano())
			logger.Warn("dotnet: " + logsanitize.ForLog(line))
			stderrOutput.WriteString(line + "\n")
		}
		wg.Done()
	}()

	logger.Sugar().Infof("executing command: %s", cmd)
	var scanner = bufio.NewScanner(stdout)
	scanner.Buffer(nil, 1024*4)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lastOutputAt.Store(time.Now().UnixNano())
		logger.Warn(logsanitize.ForLog(line))
		stdoutOutput.WriteString(line + "\n")
	}
	wg.Wait()
	close(done)
	err = cmd.Wait()
	if archivePath, archiveErr := archiveNugetRestoreBinlog(logger, solutionPath); archiveErr != nil {
		logger.Sugar().Warnf("copy NuGet restore binary log to tmp failed: src=%s err=%v", binlogPath, archiveErr)
	} else {
		logger.Sugar().Infof("NuGet restore binary log archive ready: %s", archivePath)
	}
	if err != nil {
		if idleTimeoutExceeded.Load() {
			return fmt.Errorf("dotnet restore idle timed out after %s without new output: %w\nstderr:\n%s\nstdout:\n%s",
				nugetCommandIdleTimeout,
				err,
				tailText(stderrOutput.String(), 16*1024),
				tailText(stdoutOutput.String(), 8*1024),
			)
		}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("dotnet restore timed out after %s: %w\nstderr:\n%s\nstdout:\n%s",
				time.Since(startAt).Round(time.Second),
				err,
				tailText(stderrOutput.String(), 16*1024),
				tailText(stdoutOutput.String(), 8*1024),
			)
		}
		errOutput.WriteString(fmt.Sprintf("command execution failed: %v\n", err))
		return fmt.Errorf("dotnet restore failed: %w\nstderr:\n%s\nstdout:\n%s",
			err,
			tailText(stderrOutput.String(), 16*1024),
			tailText(stdoutOutput.String(), 8*1024),
		)
	}

	return nil
}

func listNuget(ctx context.Context, task *model.InspectionTask, solutionPath string) (err error) {
	var cmdMessage string
	var stderrOutput strings.Builder
	var packageInfo ProjectPackages
	var modelVersion string
	var logger = logctx.Use(ctx)
	logNugetToolVersion(ctx, logger)
	logNugetRemoteConfigPaths(logger, solutionPath)
	if err = preflightNugetSources(ctx, logger, solutionPath); err != nil {
		return err
	}
	err = buildPackage(ctx, logger, solutionPath)
	if err != nil {
		return
	}

	listArgs := dotnetListPackageArgs(solutionPath)
	cmd := exec.CommandContext(ctx, "dotnet", listArgs...)
	cmd.Dir = filepath.Dir(solutionPath)
	if runtime.GOOS == "linux" {
		logger.Info("dotnet list sets EnableWindowsTargeting=true env for Linux compatibility")
		cmd.Env = append(os.Environ(), "EnableWindowsTargeting=true")
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("create stdout pipe failed: %w", err)
		logger.Error(err.Error())
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = fmt.Errorf("create stderr pipe failed: %w", err)
		logger.Error(err.Error())
		return
	}
	lastOutputAt := atomic.Int64{}
	lastOutputAt.Store(time.Now().UnixNano())
	idleTimeoutExceeded := atomic.Bool{}
	var stderrWg sync.WaitGroup
	stderrWg.Add(1)
	go func() {
		defer stderrWg.Done()
		scanner := bufio.NewScanner(stderr)
		scanner.Buffer(nil, 1024*4)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			lastOutputAt.Store(time.Now().UnixNano())
			logger.Debug("dotnet: " + logsanitize.ForLog(line))
			stderrOutput.WriteString(line + "\n")
		}
	}()
	logger.Sugar().Infof("executing command: %s", cmd)
	err = cmd.Start()
	if err != nil {
		// if the command is not found, we should not return error
		if errors.Is(err, exec.ErrNotFound) {
			err = _ErrDotnetNotFound
			return
		}
		err = fmt.Errorf("start command failed: %w", err)
		logger.Error(err.Error())
		return
	}
	stopWatchdog := startCommandIdleWatchdog(
		ctx, logger, cmd, "dotnet list package", solutionPath, nugetCommandIdleTimeout, &idleTimeoutExceeded, &lastOutputAt,
	)
	defer stopWatchdog()
	logger.Debug("start scanning...")
	// dotnet list --format json can be very large. Avoid per-line debug logging to reduce stream pressure.
	cmdMessage = readOutputWithHook(stdout, nil, "", func() {
		lastOutputAt.Store(time.Now().UnixNano())
	})
	waitErr := cmd.Wait()
	stderrWg.Wait()
	logger.Sugar().Infof("dotnet list stdout summary: bytes=%d lines=%d", len(cmdMessage), countLines(cmdMessage))
	if waitErr != nil {
		if idleTimeoutExceeded.Load() {
			return fmt.Errorf("dotnet list package idle timed out after %s without new output: %w\nstderr:\n%s\nstdout:\n%s",
				nugetCommandIdleTimeout,
				waitErr,
				tailText(stderrOutput.String(), 16*1024),
				tailText(cmdMessage, 8*1024),
			)
		}
		return fmt.Errorf("dotnet list package failed: %w\nstderr:\n%s\nstdout:\n%s",
			waitErr,
			tailText(stderrOutput.String(), 16*1024),
			tailText(cmdMessage, 8*1024),
		)
	}

	err = json.Unmarshal([]byte(cmdMessage), &packageInfo)
	if err != nil {
		err = fmt.Errorf("outMessage unmarshal failed: %w\nstderr:\n%s\nstdout:\n%s",
			err,
			tailText(stderrOutput.String(), 16*1024),
			tailText(cmdMessage, 8*1024),
		)
		logger.Error(err.Error())
		return
	}
	modelVersion = strconv.Itoa(packageInfo.Version)
	for _, projects := range packageInfo.Projects {
		var result []model.DependencyItem
		moduleName := filepath.Base(projects.Path)
		// 使用 map 跟踪已出现的小写包名+版本号组合，避免重复添加
		seenPackages := make(map[string]bool)

		for _, frameworks := range projects.Frameworks {
			for _, topLevelPackages := range frameworks.TopLevelPackages {
				if topLevelPackages.Id == "" {
					continue
				}
				// 将包名和版本号组合转为小写进行去重检查
				key := strings.ToLower(topLevelPackages.Id) + ":" + topLevelPackages.RequestedVersion
				if seenPackages[key] {
					// 如果已经存在相同的小写名称和版本号组合，跳过
					continue
				}
				// 标记为已出现，并添加到结果中
				seenPackages[key] = true
				result = append(result, model.DependencyItem{
					Component: model.Component{
						CompName:    topLevelPackages.Id,
						CompVersion: topLevelPackages.RequestedVersion,
						EcoRepo:     EcoRepo,
					},
					DependencyRelation: model.DependencyRelationDirect,
				})
			}
			for _, transitivePackages := range frameworks.TransitivePackages {
				if transitivePackages.Id == "" {
					continue
				}
				// 将包名和版本号组合转为小写进行去重检查
				key := strings.ToLower(transitivePackages.Id) + ":" + transitivePackages.ResolvedVersion
				if seenPackages[key] {
					// 如果已经存在相同的小写名称和版本号组合，跳过
					continue
				}
				// 标记为已出现，并添加到结果中
				seenPackages[key] = true
				result = append(result, model.DependencyItem{
					Component: model.Component{
						CompName:    transitivePackages.Id,
						CompVersion: transitivePackages.ResolvedVersion,
						EcoRepo:     EcoRepo,
					},
					DependencyRelation: model.DependencyRelationTransitive,
				})
			}
		}
		if len(result) == 0 {
			if !env.DoNotBuild {
				logger.Warn(moduleName + "::no dependencies found, backup")
			}
		}

		m := model.Module{
			ModuleName:     filepath.Base(moduleName),
			ModuleVersion:  modelVersion,
			ModulePath:     projects.Path,
			PackageManager: "nuget",
			Dependencies:   result,
		}
		task.AddModule(m)
	}
	return nil
}
