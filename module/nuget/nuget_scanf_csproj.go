package nuget

import (
	"context"
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
)

func noBuildEntrance(ctx context.Context, task *model.InspectionTask, doOld *bool) error {
	logger := logctx.Use(ctx)
	csprojpath, err := findCsproj(task.Dir())
	if err != nil {
		logger.Debug("scan csproj failed")
		return err
	}

	for _, j := range csprojpath {

		result, err := analysis(ctx, j)
		if err != nil {
			logger.Debug("analysis csproj failed")
			return err
		}

		m := model.Module{
			ModuleName:     filepath.Base(j),
			ModuleVersion:  "",
			ModulePath:     j,
			PackageManager: "nuget",
			Dependencies:   result,
		}
		*doOld = true
		task.AddModule(m)
	}
	return nil
}
func analysis(ctx context.Context, path string) (result []model.DependencyItem, err error) {
	logger := logctx.Use(ctx)
	var proj Project

	xmlData, err := os.ReadFile(path)
	if err != nil {
		logger.Debug("read file failed" + path)
		return nil, err
	}
	strings.ReplaceAll(string(xmlData), "PackageReference", "Reference")
	if err := xml.Unmarshal(xmlData, &proj); err != nil {
		logger.Debug("analysis failed")
		return nil, err
	}
	for _, pkgRef := range proj.PackageRefs {
		result = append(result, model.DependencyItem{
			Component: model.Component{
				CompName:    pkgRef.Include,
				CompVersion: pkgRef.Version,
				EcoRepo:     EcoRepo,
			},
			IsDirectDependency: true,
		})
	}
	return
}

func findCsproj(path string) ([]string, error) {
	var csprojPath []string
	return csprojPath, filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		matched, err := filepath.Match("*.csproj", d.Name())
		if err != nil {
			return err
		}
		// 检查当前路径是否是.csproj文件
		if !d.IsDir() && matched {
			csprojPath = append(csprojPath, path)
		}
		return nil
	})
}
