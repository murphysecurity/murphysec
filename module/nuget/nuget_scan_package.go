package nuget

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
)

func (this *PkgConfig) Deps() []model.DependencyItem {
	var rs []model.DependencyItem
	for _, it := range this.Package {
		if it.DevelopmentDependency {
			continue
		}
		d := model.DependencyItem{
			Component: model.Component{
				CompName:    it.Id,
				CompVersion: it.Version,
				EcoRepo:     EcoRepo,
			},
		}

		if strings.ContainsAny(d.CompVersion, "*") {
			d.CompVersion = ""
		}
		rs = append(rs, d)
	}
	return rs
}

func inspectPkgConfig(filePath string) ([]model.DependencyItem, error) {
	data, e := utils.ReadFileLimited(filePath, 4*1024*1024)
	if e != nil {
		return nil, errors.WithMessage(e, "Read packages.config failed")
	}
	var pkg PkgConfig
	if e := xml.Unmarshal(data, &pkg); e != nil {
		return nil, errors.WithMessage(e, "Parse packages.config failed")
	}
	return pkg.Deps(), nil
}

// 使用旧版规范:扫描pakcage.config
func scanPackage(task *model.InspectionTask, packagesFilePath string) error {
	dep, e := inspectPkgConfig(packagesFilePath)
	if e != nil {
		return e
	}
	m := model.Module{
		PackageManager: "nuget",
		ModuleName:     "packages.config",
		ModuleVersion:  "",
		ModulePath:     filepath.Join(task.Dir(), "packages.config"),
		Dependencies:   dep,
	}
	task.AddModule(m)
	return nil
}
