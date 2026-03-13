package nuget

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
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
	var res strings.Builder
	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(nil, 1024*1024)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
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

// 通过先运行 dotnet restore 命令，确保项目中的所有 NuGet 包依赖项被正确恢复
func buildPackage(ctx context.Context, logger *zap.Logger, solutionPath string) (err error) {
	//dotnet restore
	cmd := exec.CommandContext(ctx, "dotnet", "restore", solutionPath)
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

	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("start command failed: %v", err)
		logger.Error(err.Error())
		return
	}

	var errOutput strings.Builder
	var stderrOutput strings.Builder
	var stdoutOutput strings.Builder
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
		logger.Warn(logsanitize.ForLog(line))
		stdoutOutput.WriteString(line + "\n")
	}
	wg.Wait()
	err = cmd.Wait()
	if err != nil {
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
	err = buildPackage(ctx, logger, solutionPath)
	if err != nil {
		return
	}

	cmd := exec.CommandContext(ctx, "dotnet", "list", solutionPath, "package", "--include-transitive", "--format", "json")
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
	var stderrWg sync.WaitGroup
	stderrWg.Add(1)
	go func() {
		defer stderrWg.Done()
		scanner := bufio.NewScanner(stderr)
		scanner.Buffer(nil, 1024*4)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
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
	logger.Debug("start scanning...")
	// dotnet list --format json can be very large. Avoid per-line debug logging to reduce stream pressure.
	cmdMessage = readOutput(stdout, nil, "")
	waitErr := cmd.Wait()
	stderrWg.Wait()
	logger.Sugar().Infof("dotnet list stdout summary: bytes=%d lines=%d", len(cmdMessage), countLines(cmdMessage))
	if waitErr != nil {
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
