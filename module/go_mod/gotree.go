package go_mod

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
)

func checkNetworkEnvironment(ctx context.Context) bool {
	var goProxy []string
	if strings.Contains(os.Getenv("GOPROXY"), ",") {
		goProxy = strings.Split(os.Getenv("GOPROXY"), ",")
	} else {
		goProxy = append(goProxy, os.Getenv("GOPROXY"))
	}
	timeout := 10 * time.Second
	client := http.Client{
		Timeout: timeout, // 使用设置的超时时间
	}
	for range 3 {
		for _, j := range goProxy {
			if j == "direct" {
				continue
			}
			r, e := client.Get(j)
			if e != nil {
				logctx.Use(ctx).Warn("test network environment http get error :" + e.Error())
				continue
			}
			if r != nil && r.StatusCode == http.StatusRequestTimeout {
				logctx.Use(ctx).Warn("test network environment http get timeout :" + j)
				r.Body.Close()
				continue
			}
			logctx.Use(ctx).Debug("test network environment http get success :" + j)
			r.Body.Close()
			return true
		}
	}
	return false
}

func buildScan(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	if !checkNetworkEnvironment(ctx) {
		return errors.New("network environment error")
	}
	modFilePath := filepath.Join(task.Dir(), "go.mod")
	logger.Debug("Reading go.mod", zap.String("path", modFilePath))
	if err := goModTidy(ctx, task.Dir()); err != nil {
		logger.Error("go mod tidy error :", zap.Error(err))
		return err
	}
	modName, err := getModInfo(modFilePath)
	if err != nil {
		logger.Error("get mod info error :", zap.Error(err))
		return err
	}
	directDependencyList, err := readListCmd(ctx, task.Dir(), logger)
	if err != nil {
		logger.Error("read list cmd info error :", zap.Error(err))
		return err
	}
	nameVersionMp, indeterminacyRootList, sonTree, err := readGraphCmd(ctx, task.Dir(), directDependencyList, logger)
	if err != nil {
		logger.Error("read graph cmd info error :", zap.Error(err))
		return err
	}

	var dependencies []model.DependencyItem
	for _, j := range indeterminacyRootList {
		var packageToPackageUsed = make(map[string][]string)
		dependencie := model.DependencyItem{
			Component: model.Component{
				CompName:    j,
				CompVersion: nameVersionMp[j],
				EcoRepo:     EcoRepo,
			},
			DependencyRelation: model.DependencyRelationDirect,
		}
		logger.Debug("buildTree  start : " + j)
		dependencies = append(dependencies, buildingDependencyTree(nameVersionMp, &dependencie, sonTree, &packageToPackageUsed, logger))
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
func buildingDependencyTree(dInfo map[string]string, d *model.DependencyItem, sonTree map[string][]string, packageToPackageUsed *map[string][]string, logger *zap.Logger) model.DependencyItem {
	for name, list := range sonTree {
		if d.CompName == name {
			for _, j := range list {
				check := false
				for _, v := range (*packageToPackageUsed)[d.CompName] {
					if v == j {
						check = true
					}
				}
				if check {
					continue
				}
				mod := model.DependencyItem{
					Component: model.Component{
						CompName:    j,
						CompVersion: dInfo[j],
						EcoRepo:     EcoRepo,
					},
					DependencyRelation: model.DependencyRelationTransitive,
				}
				(*packageToPackageUsed)[d.CompName] = append((*packageToPackageUsed)[d.CompName], j)
				t := buildingDependencyTree(dInfo, &mod, sonTree, packageToPackageUsed, logger)
				t.DependencyRelation = model.DependencyRelationTransitive
				d.Dependencies = append(d.Dependencies, t)
			}
		}
	}
	return *d
}
func readListCmd(ctx context.Context, dir string, logger *zap.Logger) (map[string]bool, error) {
	var (
		err error
		cmd = exec.CommandContext(ctx, "go", "list", "-mod=readonly", "-m", "-f", "'{{if not (or .Indirect .Main)}}{{.Path}}@{{.Version}}{{end}}'", "all")
	)
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("create stdout pipe failed: %w", err)
		logger.Error(err.Error())
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = fmt.Errorf("create stderr pipe failed: %w", err)
		logger.Error(err.Error())
		return nil, err
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
			return nil, err
		}
		err = fmt.Errorf("start command failed: %w", err)
		logger.Error(err.Error())
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	var rootList = make(map[string]bool)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		text = strings.ReplaceAll(text, "'", "")
		if text == "" {
			continue
		}

		n, _, err := ParseDependencyLine(text)
		if err != nil {
			continue
		}
		rootList[n] = true
		logger.Debug("direct dependency list: " + n)
	}

	return rootList, nil
}
func readGraphCmd(ctx context.Context, dir string, directDependencyList map[string]bool, logger *zap.Logger) (map[string]string, []string, map[string][]string, error) {
	var (
		err      error
		cmd      = exec.CommandContext(ctx, "go", "mod", "graph")
		modName  string
		dInfo    = make(map[string]string)
		rootList []string
		sonTree  = make(map[string][]string)
	)
	cmd.Dir = dir
	modName, err = getModInfo(filepath.Join(dir, "go.mod"))
	if err != nil {
		logger.Error("get mod info error :", zap.Error(err))
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
			// 拿到后面的 name和version
			n, v, err := ParseDependencyLine(t[1])
			if err != nil {
				logger.Error(err.Error())
				return nil, nil, nil, err
			}

			// 如果存在  比较版本号
			// 更新最大的版本号
			if _, ok := dInfo[n]; ok {
				dInfo[n], err = comperVersion(dInfo[n], v)
				if err != nil {
					logger.Error(err.Error())
					return nil, nil, nil, err
				}
			} else {
				dInfo[n] = v
			}
			// 对比前面的字符 是不是等于包名
			// 如果是包名就是根节点
			if strings.TrimSpace(t[0]) == modName {
				if _, ok := directDependencyList[n]; ok {
					rootList = append(rootList, n)
				}
				continue
			}
			// 如果是子树级就构建子树
			name, _, err := ParseDependencyLine(t[0])
			if err != nil {
				logger.Error(err.Error())
				return nil, nil, nil, err
			}
			sonTree[name] = append(sonTree[name], n)
		}
		// logger.Debug("go: " + text)
	}
	stdout.Close()
	if err := cmd.Wait(); err != nil {
		logger.Error(err.Error())
		return nil, nil, nil, err
	}
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

func getModInfo(filepaths string) (string, error) {
	by, err := os.ReadFile(filepaths)
	if err != nil {
		return "", err
	}
	f, err := modfile.ParseLax(filepaths, by, nil)
	if err != nil {
		return "", err
	}
	return f.Module.Mod.Path, nil
}
func ParseDependencyLine(line string) (string, string, error) {
	parts := strings.Split(line, "@")
	if len(parts) != 2 {
		return "", "", errors.New("invalid dependency format: " + line)
	}
	return parts[0], parts[1], nil
}
func baseScan(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	logger := logctx.Use(ctx)
	var dependencies []model.DependencyItem
	scanMod := func(path string) error {
		var indirectMp = make(map[string]bool)
		fileTxt := ""
		modFile, err := os.Open(path)
		if err != nil {
			logger.Error("Error opening go.mod file", zap.Error(err))
			return err
		}
		defer modFile.Close()
		scanner := bufio.NewScanner(modFile)
		for scanner.Scan() {
			line := scanner.Text()
			fileTxt = fileTxt + "\n" + line
			if strings.Contains(line, "// indirect") {
				NoSpacesLine := strings.TrimSpace(line)
				nvList := strings.Split(NoSpacesLine, " ")
				indirectMp[nvList[0]] = true
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading go.mod: %v", err)
		}
		file, err := modfile.Parse(path, []byte(fileTxt), nil)
		if err != nil {
			logger.Error("Error parsing go.mod file", zap.Error(err))
			return err
		}
		modName := file.Module.Mod.Path
		for _, req := range file.Require {
			isDirectDependency := model.DependencyRelationDirect
			if _, ok := indirectMp[req.Mod.Path]; ok {
				isDirectDependency = model.DependencyRelationTransitive
			}
			dependencies = append(dependencies, model.DependencyItem{
				Component: model.Component{
					CompName:    req.Mod.Path,
					CompVersion: req.Mod.Version,
					EcoRepo:     EcoRepo,
				},
				DependencyRelation: isDirectDependency,
				IsOnline:           model.IsOnline{Value: true, Valid: true},
			})
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
	if err := scanMod(filepath.Join(task.Dir(), "go.mod")); err != nil {
		return err
	}
	return nil
}
