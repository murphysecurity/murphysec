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
	// 使用 map 跟踪已出现的小写包名+版本号组合，避免重复添加
	seenPackages := make(map[string]bool)

	for _, it := range this.Package {
		if it.DevelopmentDependency {
			continue
		}
		if it.Id == "" {
			continue
		}
		// 将包名和版本号组合转为小写进行去重检查（使用原始版本号）
		key := strings.ToLower(it.Id) + ":" + it.Version
		if seenPackages[key] {
			// 如果已经存在相同的小写名称和版本号组合，跳过
			continue
		}
		// 标记为已出现
		seenPackages[key] = true

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
