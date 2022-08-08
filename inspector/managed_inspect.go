package inspector

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/module/bundler"
	"github.com/murphysecurity/murphysec/module/cocoapods"
	"github.com/murphysecurity/murphysec/module/composer"
	"github.com/murphysecurity/murphysec/module/conan"
	"github.com/murphysecurity/murphysec/module/go_mod"
	"github.com/murphysecurity/murphysec/module/gradle"
	"github.com/murphysecurity/murphysec/module/maven"
	"github.com/murphysecurity/murphysec/module/npm"
	"github.com/murphysecurity/murphysec/module/nuget"
	"github.com/murphysecurity/murphysec/module/poetry"
	"github.com/murphysecurity/murphysec/module/python"
	"github.com/murphysecurity/murphysec/module/yarn"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

var managedInspector = []base.Inspector{
	go_mod.New(),
	maven.New(),
	npm.New(),
	gradle.New(),
	yarn.New(),
	python.New(),
	composer.New(),
	bundler.New(),
	cocoapods.New(),
	poetry.New(),
	nuget.New(),
	conan.New(),
}

type inspectorAcceptance struct {
	inspector base.Inspector
	dir       string
}

func (i inspectorAcceptance) String() string {
	return fmt.Sprintf("[%s]%s", i.inspector, i.dir)
}

func managedInspect(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	baseDir := scanTask.ProjectDir
	Logger.Info("Auto scan dir", zap.String("dir", baseDir))

	// todo: 重构，随着检查器越来越多，这里越来越慢
	var inspectorAcceptances []inspectorAcceptance
	for _, inspector := range managedInspector {
		e := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				Logger.Warn("item is nil", zap.Error(err))
				return nil
			}
			if !d.IsDir() {
				return nil
			}
			if d.IsDir() && dirShouldIgnore(d.Name()) {
				return filepath.SkipDir
			}
			if relDir, e := filepath.Rel(baseDir, path); e == nil {
				if strings.Count(filepath.ToSlash(relDir), "/") > 5 {
					return filepath.SkipDir
				}
			} else {
				return nil
			}
			if inspector.CheckDir(path) {
				inspectorAcceptances = append(inspectorAcceptances, inspectorAcceptance{inspector, path})
				if !inspector.SupportFeature(base.FeatureAllowNested) {
					return filepath.SkipDir
				}
				return nil
			}
			return nil
		})
		if e != nil {
			return e
		}
	}
	Logger.Sugar().Infof("Found %d directories, in %v", len(inspectorAcceptances), time.Now().Sub(scanTask.StartTime))
	for idx, acceptance := range inspectorAcceptances {
		st := time.Now()
		c := model.WithInspectorTask(ctx, acceptance.dir)
		c = utils.WithLogger(c, Logger.Named(fmt.Sprintf("%s-%d", acceptance.inspector.String(), idx)))
		e := acceptance.inspector.InspectProject(c)
		Logger.Sugar().Infof("%v, duration: %v", acceptance, time.Now().Sub(st))
		if e != nil {
			Logger.Error("InspectError", zap.Error(e))
		}
	}
	return nil
}
