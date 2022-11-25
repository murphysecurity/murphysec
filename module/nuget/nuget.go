package nuget

import (
	"context"
	"encoding/xml"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
)

type Inspector struct{}

func (i *Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

func (i *Inspector) String() string {
	return "Nuget"
}
func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "packages.config"))
}
func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectionTask(ctx)
	dep, e := inspectPkgConfig(filepath.Join(task.Dir(), "packages.config"))
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

type PkgConfig struct {
	XMLName xml.Name `xml:"packages"`
	Package []struct {
		Id                    string `xml:"id,attr"`
		Version               string `xml:"version,attr"`
		DevelopmentDependency bool   `xml:"developmentDependency,attr"`
	} `xml:"package"`
}

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

var EcoRepo = model.EcoRepo{
	Ecosystem:  "nuget",
	Repository: "",
}
