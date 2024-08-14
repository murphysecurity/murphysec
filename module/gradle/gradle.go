package gradle

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/sl"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/repeale/fp-go"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Inspector struct{}

func (Inspector) SupportFeature(_ model.InspectorFeature) bool {
	return false
}

func (Inspector) String() string {
	return "Gradle"
}

func (Inspector) InspectProject(ctx context.Context) error {
	var logger = logctx.Use(ctx).Sugar()
	var rs []model.Module
	task := model.UseInspectionTask(ctx)
	dir := task.Dir()
	logger.Debugf("gradle inspect dir: %s", dir)
	useGradle := true
	gradleEnv, e := DetectGradleEnv(ctx, dir)
	if e != nil {
		logger.Infof("check gradle failed: %s", e.Error())
		logger.Warnf("Gradle disabled")
		useGradle = false
	}
	if useGradle {
		logger.Info(gradleEnv.Version.String())
		rs, e = evalGradleDependencies(ctx, dir, gradleEnv)
	}
	if len(rs) == 0 && !env.ScannerScan {
		// if no module find, use backup solution
		if m := backupParser(ctx, dir); m != nil {
			tm := m.BaseModule(dir)
			tm.ScanStrategy = model.ScanStrategyBackup
			rs = append(rs, tm)
		}
	}
	env.ScannerShouldEnableGradleBackupScan = true
	for _, i := range rs {
		if len(i.Dependencies) != 0 {
			env.ScannerShouldEnableGradleBackupScan = false
			break
		}
	}
	for _, it := range rs {
		task.AddModule(it)
	}
	return nil
}

func gradleProjectName2PathComponent(baseDir string, projectName string) string {
	part2 := strings.Join(fp.Filter(sl.StringNotEmpty)(strings.Split(projectName, ":")), "/")
	if f := filepath.Join(baseDir, part2, "build.gradle.kts"); utils.IsFile(f) {
		return f
	}
	return filepath.Join(baseDir, part2, "build.gradle")
}

func backupParser(ctx context.Context, dir string) *GradleDependencyInfo {
	var logger = logctx.Use(ctx).Sugar()
	var dep []DepElement
	e := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || !d.Type().IsRegular() || !strings.HasPrefix(d.Name(), "build.gradle") {
			return nil
		}
		if d.Name() == "build.gradle.kts" {
			data, e := os.ReadFile(path)
			if e != nil {
				logger.Errorf("Read gradle file failed: %s", e.Error())
				return nil
			}
			dep = append(dep, parseGradleKts(string(data))...)
			return nil
		}
		if d.Name() == "build.gradle" {
			data, e := os.ReadFile(path)
			if e != nil {
				logger.Errorf("Read gradle file failed: %s", e.Error())
				return nil
			}
			dep = append(dep, parseGradleGroovy(string(data))...)
			return nil
		}
		return nil
	})
	if e != nil {
		logger.Warnf("Walk: %v", e)
	}
	if len(dep) != 0 {
		return &GradleDependencyInfo{
			ProjectName:  fmt.Sprintf("GradleProject-%s", filepath.Base(dir)),
			Dependencies: dep,
		}
	}
	return nil
}

var gradleBuildFiles = []string{"build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"}

func (Inspector) CheckDir(dir string) bool {
	for _, it := range gradleBuildFiles {
		info, e := os.Stat(filepath.Join(dir, it))
		if e == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

var implTaskNamePattern = regexp.MustCompile(`(?i)(?:release|debug|)(?:kotlin)?(?:compile|runtime)Classpath`)

//goland:noinspection GoNameStartsWithPackageName
type GradleDependencyInfo struct {
	ProjectName  string       `json:"project_name"`
	Dependencies []DepElement `json:"dependencies,omitempty"`
}

func (g *GradleDependencyInfo) BaseModule(basePath string) model.Module {
	return model.Module{
		PackageManager: "gradle",
		ModuleName:     g.ProjectName,
		Dependencies:   convDep(g.Dependencies),
		ModulePath:     gradleProjectName2PathComponent(basePath, g.ProjectName),
	}
}

func convDep(input []DepElement) []model.DependencyItem {
	var r = _convDep(input)
	for i := range r {
		r[i].IsDirectDependency = true
	}
	return r
}

func _convDep(input []DepElement) []model.DependencyItem {
	var rs []model.DependencyItem
	for _, it := range input {
		rs = append(rs, model.DependencyItem{
			Component: model.Component{
				CompName:    it.CompName(),
				CompVersion: it.Version,
				EcoRepo:     EcoRepo,
			},
			Dependencies: _convDep(it.Children),
		})
	}
	return rs
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "maven",
	Repository: "",
}

type DepElement struct {
	GroupId    string       `json:"group_id"`
	ArtifactId string       `json:"artifact_id"`
	Version    string       `json:"version"`
	Children   []DepElement `json:"children,omitempty"`
}

func (d DepElement) CompName() string {
	return fmt.Sprintf("%s:%s", d.GroupId, d.ArtifactId)
}

//go:embed print-dep.gradle
var scriptPrintDep string

func evalGradleDependencies(ctx context.Context, dir string, info *GradleEnv) (modules []model.Module, e error) {
	var logger = logctx.Use(ctx).Sugar()
	defer func() {
		if e != nil {
			e = fmt.Errorf("generateGradleDepFiles: %w", e)
		}
	}()
	td, e := os.MkdirTemp("", "mpsgradle-")
	if e != nil {
		e = fmt.Errorf("failed to create temp directory: %w", e)
		return
	}
	var tdf = filepath.Join(td, "print-dep.gradle")
	e = os.WriteFile(tdf, []byte(scriptPrintDep), 0o644)
	if e != nil {
		e = fmt.Errorf("failed to write temp file: %w", e)
		return
	}
	logger.Infof("temp file written, %s", tdf)

	var cmd = info.ExecuteContext(ctx, "-I", tdf, "generateDependencyFile", "--info")
	cmd.Dir = dir
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	stdout, stderr, e = utils.ExecGetStdOutErr(ctx, cmd)
	if e != nil {
		return
	}
	e = cmd.Start()
	if e != nil {
		e = fmt.Errorf("failed to start gradle: %w", e)
		return
	}
	var wg sync.WaitGroup
	defer wg.Wait()
	forwardStdAsync(ctx, "gradle", stdout, &wg)
	forwardStdAsync(ctx, "gradle[E]", stderr, &wg)
	e = cmd.Wait()
	if e != nil {
		e = fmt.Errorf("gradle failed: %w", e)
		logger.Desugar().Error(e.Error())
		return
	}
	logger.Info("collecting results...")
	var handler, collector = parseGradleScriptOutputAsyncBuilder(ctx)
	e = walkGradleScriptOutput(ctx, dir, handler)
	if e != nil {
		return
	}
	modules, e = collector()
	if e != nil {
		e = fmt.Errorf("collecting results failed: %w", e)
		return
	}
	return
}

func forwardStdAsync(ctx context.Context, prefix string, reader io.ReadCloser, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = reader.Close() }()
		var logger = logctx.Use(ctx).Sugar()
		var scanner = bufio.NewScanner(reader)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			if scanner.Err() != nil {
				logger.Errorf("%s: %v", prefix, scanner.Err())
				break
			}
			logger.Debugf("%s: %s", prefix, scanner.Text())
		}
	}()
}

func walkGradleScriptOutput(ctx context.Context, dir string, handler func(dir string) error) error {
	var logger = logctx.Use(ctx).Sugar()
	var e = filepath.WalkDir(dir, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") || d.Name() == "__MACOSX" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != "dependency-tree-mp.yaml" {
			return nil
		}
		logger.Debugf("found dependency-tree-mp.yaml: %s", path)
		e = handler(path)
		return nil
	})
	if e != nil {
		e = fmt.Errorf("walkGradleScriptOutput: %w", e)
	}
	return e
}

func parseGradleScriptOutputAsyncBuilder(ctx context.Context) (handler func(path string) error, collector func() ([]model.Module, error)) {
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(16)
	var mutex sync.Mutex
	var modules []model.Module
	handler = func(path string) error {
		eg.Go(func() (e error) {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			f, e := os.Open(path)
			if e != nil {
				return e
			}
			defer func() { _ = f.Close() }()
			var _modules []model.Module
			_modules, e = decodeGradleScriptOutput(f, path)
			if e != nil {
				return
			}
			mutex.Lock()
			modules = append(modules, _modules...)
			mutex.Unlock()
			return nil
		})
		return nil
	}
	return handler, func() ([]model.Module, error) {
		return modules, eg.Wait()
	}
}

func decodeGradleScriptOutput(reader io.Reader, dir string) (modules []model.Module, e error) {
	var decoder = yaml.NewDecoder(reader)
	var data dtoProjectData
	e = decoder.Decode(&data)
	if e != nil {
		return
	}
	for _, configuration := range data.Configurations {
		var online = model.IsOnlineFalse()
		if !strings.Contains(strings.ToLower(configuration.Configuration), "test") {
			online = model.IsOnlineTrue()
		}
		var module = model.Module{
			PackageManager: "gradle",
			ModuleName:     data.Project + ":" + configuration.Configuration,
			Dependencies:   fp.Map(func(it dtoItem) model.DependencyItem { return it.toItem(online) })(configuration.Dependencies),
			ScanStrategy:   model.ScanStrategyNormal,
			ModulePath:     path.Join(filepath.Join(dir, "../../build.gradle"), ":"+configuration.Configuration),
		}
		modules = append(modules, module)
	}
	return
}

type dtoProjectData struct {
	Project        string             `yaml:"project" json:"project"`
	Configurations []dtoConfiguration `yaml:"configurations" json:"configurations"`
}

type dtoConfiguration struct {
	Configuration string    `yaml:"configuration" json:"configuration"`
	Dependencies  []dtoItem `yaml:"dependencies" json:"dependencies"`
}

type dtoItem struct {
	Group    string    `yaml:"group" json:"group"`
	Name     string    `yaml:"name" json:"name"`
	Version  string    `yaml:"version" json:"version"`
	Children []dtoItem `yaml:"children" json:"children"`
}

func (d dtoItem) toItem(online model.IsOnline) model.DependencyItem {
	return model.DependencyItem{
		Component: model.Component{
			CompName:    d.Group + ":" + d.Name,
			CompVersion: d.Version,
			EcoRepo:     EcoRepo,
		},
		Dependencies: fp.Map(func(t dtoItem) model.DependencyItem { return t.toItem(online) })(d.Children),
		IsOnline:     online,
	}
}
