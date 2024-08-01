package ivy

import (
	"bufio"
	"context"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	file, e := os.Open(ivyPath)
	if e != nil {
		return fmt.Errorf("open file: %w", e)
	}
	defer func() { _ = file.Close() }()
	module, e := readIvyXml(ctx, bufio.NewReader(file))
	if e != nil {
		return fmt.Errorf("read ivy.xml: %w", e)
	}
	module.ModulePath = ivyPath
	task.AddModule(*module)
	return nil
}

func readIvyXml(ctx context.Context, reader io.Reader) (*model.Module, error) {
	root, e := xmlquery.Parse(reader)
	if e != nil {
		return nil, fmt.Errorf("parse xml: %w", e)
	}
	module := model.Module{
		PackageManager: "ivy",
		ModuleName:     "<NoName>",
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
		version := node.SelectAttr("rev")
		if org == "" || name == "" {
			return
		}
		// if version is not specified literally, leave it empty
		if strings.Contains(version, "$") {
			version = ""
		}
		module.Dependencies = append(module.Dependencies, model.DependencyItem{
			Component: model.Component{
				CompName:    org + ":" + name,
				CompVersion: version,
				EcoRepo:     EcoRepo,
			},
			IsOnline:           model.IsOnlineTrue(),
			IsDirectDependency: true,
		})
	})
	return &module, nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "maven",
	Repository: "",
}
