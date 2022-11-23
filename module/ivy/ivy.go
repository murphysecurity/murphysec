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
	task := model.UseInspectionTask(ctx)
	ivyPath := filepath.Join(task.Dir(), "ivy.xml")
	data, e := os.ReadFile(ivyPath)
	if e != nil {
		return fmt.Errorf("open ivy file: %w", e)
	}
	root, e := xmlquery.Parse(bytes.NewReader(data))
	if e != nil {
		return fmt.Errorf("parse xml: %w", e)
	}
	module := model.Module{
		PackageManager: "ivy",
		ModuleName:     "<NoName>",
		ModulePath:     ivyPath,
		Dependencies:   make([]model.DependencyItem, 0),
		ScanStrategy:   model.ScanStrategyBackup,
	}

	if infoNode := xmlquery.FindOne(root, "//ivy-module/info"); infoNode != nil {
		org := infoNode.SelectAttr("organisation")
		name := infoNode.SelectAttr("module")
		ver := infoNode.SelectAttr("revision")
		module.ModuleName = org + ":" + name
		module.ModuleVersion = ver
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
		module.Dependencies = append(module.Dependencies, model.DependencyItem{
			Component: model.Component{
				CompName:    org + ":" + name,
				CompVersion: version,
				EcoRepo:     EcoRepo,
			},
		})
	})
	task.AddModule(module)
	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "maven",
	Repository: "",
}
