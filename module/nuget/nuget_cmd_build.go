package nuget

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	"github.com/murphysecurity/murphysec/model"
	"go.uber.org/zap"
)

var _ErrDotnetNotFound = errors.New("dotnet not found")

func multipleBuilds(ctx context.Context, task *model.InspectionTask) error {
	logger := logctx.Use(ctx)
	filePath, err := findCLNList(task.Dir())
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	numCPU := utils.Coerce(runtime.NumCPU(), 1, 4)
	var wg sync.WaitGroup
	ch := make(chan string, len(filePath))
	for _, j := range filePath {
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
				}
			}
		}()
	}
	wg.Wait()
	return nil

}
func buildEntrance(ctx context.Context, task *model.InspectionTask, directory string) error {
	logger := logctx.Use(ctx)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	e := listNuget(ctx, task, directory)
	if e != nil {
		if errors.Is(e, _ErrDotnetNotFound) {
			logger.Warn("Dotnet not found, skip DotnetList")
			return e
		} else {
			// log it and go on
			logger.Warn("Dotnet list failed"+directory, zap.Error(e))
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
func readOutput(pipe io.ReadCloser) string {
	var res string
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		res += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("reading output failed: %v\n", err)
	}
	return res
}

// 通过先运行 dotnet restore 命令，确保项目中的所有 NuGet 包依赖项被正确恢复
func buildPackage(ctx context.Context, logger *zap.Logger, directory string) (err error) {
	//dotnet restore
	cmd := exec.CommandContext(ctx, "dotnet", "restore")
	cmd.Dir = directory
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Buffer(nil, 1024*4)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logger.Warn("dotnet: " + scanner.Text())
		}
		wg.Done()
	}()

	logger.Sugar().Infof("executing command: %s", cmd)
	var scanner = bufio.NewScanner(stdout)
	scanner.Buffer(nil, 1024*4)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		logger.Warn(scanner.Text())
	}
	wg.Wait()
	err = cmd.Wait()
	if err != nil {
		errOutput.WriteString(fmt.Sprintf("command execution failed: %v\n", err))
		return err
	}

	return nil
}

func listNuget(ctx context.Context, task *model.InspectionTask, directory string) (err error) {
	dir := directory
	var cmdMessage string
	var packageInfo ProjectPackages
	var modelVersion string
	var logger = logctx.Use(ctx)
	err = buildPackage(ctx, logger, dir)
	if err != nil {
		return
	}

	cmd := exec.CommandContext(ctx, "dotnet", "list", "package", "--include-transitive", "--format", "json")
	cmd.Dir = dir
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
	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Buffer(nil, 1024*4)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logger.Warn("dotnet: " + scanner.Text())
		}
	}()
	logger.Sugar().Infof("executing command: %s", cmd)
	var scanner = bufio.NewScanner(stdout)
	scanner.Buffer(nil, 1024*4)
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
	cmdMessage = readOutput(stdout)

	err = json.Unmarshal([]byte(cmdMessage), &packageInfo)
	if err != nil {
		err = fmt.Errorf("outMessage unmarshal failed: %w", err)
		logger.Error(err.Error())
		return
	}
	modelVersion = strconv.Itoa(packageInfo.Version)
	for _, projects := range packageInfo.Projects {
		var result []model.DependencyItem
		moduleName := filepath.Base(projects.Path)
		for _, frameworks := range projects.Frameworks {
			for _, topLevelPackages := range frameworks.TopLevelPackages {
				result = append(result, model.DependencyItem{
					Component: model.Component{
						CompName:    topLevelPackages.Id,
						CompVersion: topLevelPackages.RequestedVersion,
						EcoRepo:     EcoRepo,
					},
					IsDirectDependency: true,
				})
			}
			for _, transitivePackages := range frameworks.TransitivePackages {
				result = append(result, model.DependencyItem{
					Component: model.Component{
						CompName:    transitivePackages.Id,
						CompVersion: transitivePackages.ResolvedVersion,
						EcoRepo:     EcoRepo,
					},
					IsDirectDependency: false,
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
	_ = cmd.Wait()
	return nil
}
