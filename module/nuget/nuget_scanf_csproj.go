package nuget

import (
	"context"
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

var versionPattern = regexp.MustCompile(`Version=([^,]+)(?:,|$)`)

func analysis(ctx context.Context, path string) (result []model.DependencyItem, err error) {
	logger := logctx.Use(ctx)
	var proj Project

	xmlData, err := os.ReadFile(path)
	if err != nil {
		logger.Debug("read file failed" + path)
		return nil, err
	}

	if err := xml.Unmarshal(xmlData, &proj); err != nil {
		logger.Debug("analysis failed")
		return nil, err
	}

	for _, j := range proj.Reference {

		var mod struct {
			Include string `xml:"Include,attr"`
			Version string `xml:"Version,attr"`
		}

		var includePackage string
		if index := strings.Index(j.Include, ","); index != -1 && !strings.Contains(j.Include[:1], " ") {
			includePackage = j.Include[:index]

		} else {
			continue
		}
		versionMatches := versionPattern.FindStringSubmatch(j.Include)
		if len(versionMatches) == 0 {
			continue
		}
		logger.Error(includePackage)
		mod.Include = includePackage
		mod.Version = versionMatches[1]
		proj.PackageRefs = append(proj.PackageRefs, mod)

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
