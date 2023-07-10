package npm

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Npm"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "package-lock.json"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
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
	dir := model.UseInspectionTask(ctx).Dir()
	packagePath := filepath.Join(dir, "package.json")
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

	lockfilePath := filepath.Join(dir, "package-lock.json")
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
		module.ModuleName = parsed.Name
		module.ModuleVersion = parsed.Version
		module.Dependencies = parsed.Deps
		return []model.Module{module}, nil
	}

	module.ModuleName = packageFile.Name
	module.ModuleVersion = packageFile.Version
	var requires = lo.Keys(packageFile.Dependencies)
	requires = append(requires, lo.Keys(packageFile.DevDependencies)...)
	requires = lo.Uniq(requires)
	deps, e := processV1Lockfile(data, requires)
	if e != nil {
		return nil, e
	}
	module.Dependencies = utils.NoNilSlice(deps)
	return []model.Module{module}, nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}
