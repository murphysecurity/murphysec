package npm

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
)

type Inspector struct{}

const PackageFileName = "package.json"
const LockFileName = "package-lock.json"

var autoBuildDisabled = errors.New("auto build disabled")

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

func (Inspector) String() string {
	return "Npm"
}

func (Inspector) CheckDir(dir string) bool {
	if utils.IsFile(filepath.Join(dir, LockFileName)) {
		return true
	}
	if !utils.IsFile(filepath.Join(dir, PackageFileName)) {
		return false
	}
	if utils.IsFile(filepath.Join(dir, "yarn.lock")) {
		return false
	}
	if utils.IsFile(filepath.Join(dir, "pnpm-lock.yaml")) {
		return false
	}
	return true
}

func (Inspector) InspectProject(ctx context.Context) error {
	m, e := ScanNpmProject(ctx)
	if e != nil {
		return e
	}
	for _, it := range m {
		model.UseInspectionTask(ctx).AddModule(it)
	}
	return nil
}

func ScanNpmProject(ctx context.Context) ([]model.Module, error) {
	logger := logctx.Use(ctx)
	dir := model.UseInspectionTask(ctx).Dir()
	packagePath := filepath.Join(dir, PackageFileName)
	module := model.Module{
		PackageManager: "npm",
		ModuleName:     "",
		ModuleVersion:  "",
		ModulePath:     packagePath,
	}

	data, e := os.ReadFile(packagePath)
	if e != nil {
		return nil, fmt.Errorf("reading package file: %w", e)
	}
	packageFile, e := parsePkgFile(data)
	if e != nil {
		return nil, e
	}
	module.ModuleName = packageFile.Name
	module.ModuleVersion = packageFile.Version
	lockfilePath := filepath.Join(dir, LockFileName)
	if !utils.IsPathExist(lockfilePath) {
		if env.DoNotBuild {
			logger.Info("lockfile not found, and auto build is disabled, skip")
			return make([]model.Module, 0), nil
		}
		e = doNpmInstallInDir(ctx, dir)
		if errors.Is(e, autoBuildDisabled) {
			logger.Warn("fallback to read package.json only")
			for k, v := range packageFile.Dependencies {
				module.Dependencies = append(module.Dependencies, model.DependencyItem{
					Component: model.Component{
						CompName:    k,
						CompVersion: v,
						EcoRepo:     EcoRepo,
					},
					IsDirectDependency: true,
					IsOnline:           model.IsOnlineTrue(),
				})
			}
			for k, v := range packageFile.DevDependencies {
				module.Dependencies = append(module.Dependencies, model.DependencyItem{
					Component: model.Component{
						CompName:    k,
						CompVersion: v,
						EcoRepo:     EcoRepo,
					},
					IsDirectDependency: true,
					IsOnline:           model.IsOnlineFalse(),
				})
			}
			return []model.Module{module}, nil
		} else if e != nil {
			logger.Warn("npm install failed, skip")
			return make([]model.Module, 0), nil
		}
	}
	data, e = os.ReadFile(lockfilePath)
	if e != nil {
		return nil, fmt.Errorf("reading package-lock file: %w", e)
	}
	lockfileVer, e := parseLockfileVersion(data)
	if e != nil {
		return nil, e
	}
	if lockfileVer == 3 {
		parsed, e := processLockfileV3(data)
		if e != nil {
			return nil, fmt.Errorf("v3lockfile: %w", e)
		}
		module.ModuleName = parsed.CompName
		module.ModuleVersion = parsed.CompVersion
		module.Dependencies = parsed.Dependencies
		return []model.Module{module}, nil
	}

	deps, e := processV1Lockfile(data, packageFile)
	if e != nil {
		return nil, e
	}
	module.Dependencies = utils.NoNilSlice(deps)
	return []model.Module{module}, nil
}

func doNpmInstallInDir(ctx context.Context, dir string) error {
	logger := logctx.Use(ctx)
	if env.DoNotBuild {
		logger.Warn("lockfile not found, and auto build is disabled, skip")
		return autoBuildDisabled
	}
	logger.Info("lockfile not found, do npm install...")
	cmd := exec.CommandContext(ctx, "npm", "i", "--package-lock-only")
	cmd.Dir = dir
	logger.Info("command: npm i --package-lock-only")
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		panic(e)
	}
	go func() {
		logger.Debug("npm log forwarding begin...")
		defer func() { logger.Debug("npm log forwarding end.") }()
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 1024), 1024)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			logger.Debug("npm:", zap.String("log", scanner.Text()))
		}
	}()
	logger.Debug("command start")
	e = cmd.Start()
	if e != nil {
		return fmt.Errorf("start npm install command failed: %w", e)
	}
	logger.Debug("waiting...")
	e = cmd.Wait()
	if e != nil {
		return fmt.Errorf("wait npm install command terminate: %w", e)
	}
	logger.Debug("done.")
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}
