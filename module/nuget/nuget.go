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
	if utils.IsFile(filepath.Join(dir, "packages.config")) {
		return true
	}
	return false
}
func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	dep, e := inspectPkgConfig(filepath.Join(task.ScanDir, "packages.config"))
	if e != nil {
		return e
	}
	m := model.Module{
		PackageManager: model.PmNuget,
		Language:       model.DotNet,
		Name:           "packages.config",
		Version:        "",
		RelativePath:   filepath.Join(task.ProjectDir, "packages.config"),
		Dependencies:   dep,
		RuntimeInfo:    nil,
	}
	task.AddModule(m)
	return nil
}

func (i *Inspector) PackageManagerType() model.PackageManagerType {
	return model.PmNuget
}

type PkgConfig struct {
	XMLName xml.Name `xml:"packages"`
	Package []struct {
		Id                    string `xml:"id,attr"`
		Version               string `xml:"version,attr"`
		DevelopmentDependency bool   `xml:"developmentDependency,attr"`
	} `xml:"package"`
}

func (this *PkgConfig) Deps() []model.Dependency {
	rs := make([]model.Dependency, 0)
	for _, it := range this.Package {
		if it.DevelopmentDependency {
			continue
		}
		d := model.Dependency{
			Name:    it.Id,
			Version: it.Version,
		}
		if strings.ContainsAny(d.Version, "*") {
			d.Version = ""
		}
		rs = append(rs, d)
	}
	return rs
}

func inspectPkgConfig(filePath string) ([]model.Dependency, error) {
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
