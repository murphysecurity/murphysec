package ivy

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"os"
	"path/filepath"
)

type Inspector struct{}

func (Inspector) String() string {
	return "Ivy"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "ivy.xml"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	ivyPath := filepath.Join(task.ScanDir, "ivy.xml")
	data, e := os.ReadFile(ivyPath)
	if e != nil {
		return fmt.Errorf("open ivy file: %w", e)
	}
	root, e := xmlquery.Parse(bytes.NewReader(data))
	if e != nil {
		return fmt.Errorf("parse xml: %w", e)
	}
	module := model.Module{
		PackageManager: model.PmIvy,
		Language:       model.Java,
		Name:           "<NoName>",
		Version:        "",
		RelativePath:   ivyPath,
		Dependencies:   make([]model.Dependency, 0),
		ScanStrategy:   model.ScanStrategyBackup,
	}

	if infoNode := xmlquery.FindOne(root, "//ivy-module/info"); infoNode != nil {
		org := infoNode.SelectAttr("organisation")
		name := infoNode.SelectAttr("module")
		ver := infoNode.SelectAttr("revision")
		module.Name = org + ":" + name
		module.Version = ver
	}
	xmlquery.FindEach(root, "//ivy-module/dependencies/dependency", func(i int, node *xmlquery.Node) {
		//var org, name, version string
		org := node.SelectAttr("organisation")
		if org == "" {
			org = node.SelectAttr("org")
		}
		name := node.SelectAttr("name")
		version := node.SelectAttr("version")
		if org == "" || name == "" {
			return
		}
		module.Dependencies = append(module.Dependencies, model.Dependency{
			Name:    org + ":" + name,
			Version: version,
		})
	})
	task.AddModule(module)
	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}
