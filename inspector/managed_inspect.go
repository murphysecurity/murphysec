package inspector

import (
	"context"
	_ "embed"
	"fmt"
	"io/fs"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/bundler"
	"murphysec-cli-simple/module/cocoapods"
	"murphysec-cli-simple/module/composer"
	"murphysec-cli-simple/module/go_mod"
	"murphysec-cli-simple/module/gradle"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/module/npm"
	"murphysec-cli-simple/module/poetry"
	"murphysec-cli-simple/module/python"
	"murphysec-cli-simple/module/yarn"
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
	logger.Info.Println("Auto scan dir:", baseDir)

	// todo: 重构，随着检查器越来越多，这里越来越慢
	var inspectorAcceptances []inspectorAcceptance
	for _, inspector := range managedInspector {
		e := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				logger.Warn.Println(fmt.Sprintf("Can't walk into: %s due an error: %s", path, err.Error()))
				return nil
			}
			if d == nil {
				return nil
			}
			if d.IsDir() && dirIgnored(d.Name()) {
				return filepath.SkipDir
			}
			if relDir, e := filepath.Rel(baseDir, path); e == nil {
				if strings.Count(filepath.ToSlash(relDir), "/") > 3 {
					return filepath.SkipDir
				}
			} else {
				return nil
			}
			if inspector.CheckDir(path) {
				inspectorAcceptances = append(inspectorAcceptances, inspectorAcceptance{inspector, path})
				return filepath.SkipDir
			}
			return nil
		})
		if e != nil {
			return e
		}
	}

	logger.Info.Printf("Found %d directories, in %v", len(inspectorAcceptances), time.Now().Sub(scanTask.StartTime))
	for _, it := range inspectorAcceptances {
		logger.Debug.Println(it)
	}
	for _, acceptance := range inspectorAcceptances {
		st := time.Now()
		c := model.WithInspectorTask(ctx, acceptance.dir)
		e := acceptance.inspector.InspectProject(c)
		logger.Info.Printf("%v, duration: %v", acceptance, time.Now().Sub(st))
		if e != nil {
			logger.Err.Println("InspectorError:", e.Error())
		}
	}
	return nil
}
