package go_mod

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
)

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	modFilePath := filepath.Join(task.Dir(), "go.mod")
	logger.Debug("Reading go.mod", zap.String("path", modFilePath))
	modName, _, err := getModInfo(modFilePath)
	if err != nil {
		return err
	}
	nameVersionMp, rootList, sonTree, err := readCmd(ctx, task.Dir(), logger)
	if err != nil {
		return err
	}
	var dependencies []model.DependencyItem
	var packageToPackageUsed = make(map[string]string)
	for _, j := range rootList {
		dependencie := model.DependencyItem{
			Component: model.Component{
				CompName:    j,
				CompVersion: nameVersionMp[j],
				EcoRepo:     EcoRepo,
			},
			IsDirectDependency: true,
		}
		logger.Debug("buildTree  start : " + j)
		buildingDependencyTree(nameVersionMp, &dependencie, sonTree, packageToPackageUsed, logger)

		dependencies = append(dependencies, dependencie)
	}
	m := model.Module{
		PackageManager: "gomod",
		ModulePath:     filepath.Join(task.Dir(), "go.mod"),
		ModuleName:     modName,
		Dependencies:   dependencies,
	}

	task.AddModule(m)
	return nil
}
func buildingDependencyTree(dInfo map[string]string, d *model.DependencyItem, sonTree map[string][]string, packageToPackageUsed map[string]string, logger *zap.Logger) {
	for name, list := range sonTree {
		if d.CompName == name {
			for _, j := range list {
				if n, ok := packageToPackageUsed[j]; ok && n == d.CompName {
					continue
				}
				mod := model.DependencyItem{
					Component: model.Component{
						CompName:    j,
						CompVersion: dInfo[j],
					},
					IsDirectDependency: false,
				}
				d.Dependencies = append(d.Dependencies, mod)
				packageToPackageUsed[j] = d.CompName
				packageToPackageUsed[d.CompName] = j
				buildingDependencyTree(dInfo, &mod, sonTree, packageToPackageUsed, logger)
			}
		}
	}
}
func readCmd(ctx context.Context, dir string, logger *zap.Logger) (map[string]string, []string, map[string][]string, error) {
	var (
		err      error
		cmd      = exec.CommandContext(ctx, "go", "mod", "graph")
		modName  string
		dInfo    = make(map[string]string)
		rootList []string
		sonTree  = make(map[string][]string)
		isUsed   = make(map[string]bool)
	)
	cmd.Dir = dir
	modName, _, err = getModInfo(dir)
	if err != nil {
		return nil, nil, nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("create stdout pipe failed: %w", err)
		logger.Error(err.Error())
		return nil, nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = fmt.Errorf("create stderr pipe failed: %w", err)
		logger.Error(err.Error())
		return nil, nil, nil, err
	}

	go func() {
		defer stderr.Close()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logger.Warn("go: " + scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		// if the command is not found, we should not return error
		if errors.Is(err, exec.ErrNotFound) {
			err = _ErrGoNotFound
			return nil, nil, nil, err
		}
		err = fmt.Errorf("start command failed: %w", err)
		logger.Error(err.Error())
		return nil, nil, nil, err
	}
	scanner := bufio.NewScanner(stdout)
	var lineNumber int = 0
	for scanner.Scan() {
		lineNumber++
		text := scanner.Text()
		if text == "" {
			continue
		}
		t := strings.Split(text, " ")
		if len(t) == 2 {
			//拿到后面的 name和version
			n, v, err := ParseDependencyLine(t[1])
			if err != nil {
				logger.Error(err.Error())
				return nil, nil, nil, err
			}

			//如果存在  比较版本号
			//更新最大的版本号
			if _, ok := dInfo[n]; ok {
				dInfo[n], err = comperVersion(dInfo[n], v)
				if err != nil {
					logger.Error(err.Error())
					return nil, nil, nil, err
				}
			} else {
				dInfo[n] = v
			}
			//对比前面的字符 是不是等于包名
			//如果是包名就是根节点
			if strings.TrimSpace(t[0]) == modName {
				rootList = append(rootList, n)
				continue
			}
			//如果是子树级就构建子树
			//避免 环
			if _, ok := isUsed[n]; !ok {
				name, _, err := ParseDependencyLine(t[0])
				if err != nil {
					return nil, nil, nil, err
				}
				isUsed[n] = true
				sonTree[name] = append(sonTree[name], n)
			}

		}
		logger.Debug("go: " + text)
	}
	stdout.Close()
	cmd.Wait()
	return dInfo, rootList, sonTree, nil
}

func comperVersion(version1, version2 string) (string, error) {
	// Parse the version strings into semver.Version objects
	v1, err := semver.NewVersion(version1)
	if err != nil {

		return "", err
	}

	v2, err := semver.NewVersion(version2)
	if err != nil {

		return "", err
	}

	// Compare the two versions
	if !v1.LessThan(v2) {
		return version2, nil
	}
	return version2, nil
}
func getModInfo(filepath string) (string, string, error) {
	by, _ := os.ReadFile("./go.mod")
	f, _ := modfile.Parse("./go.mod", by, nil)
	return f.Module.Mod.Path, f.Go.Version, nil
}
func ParseDependencyLine(line string) (string, string, error) {
	parts := strings.Split(line, "@")
	if len(parts) != 2 {
		return "", "", errors.New("invalid dependency format: " + line)
	}
	return parts[0], parts[1], nil
}
