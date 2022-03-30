package inspector

import (
	"io/fs"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/go_mod"
	"murphysec-cli-simple/module/gradle"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/module/npm"
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
}

// 受管理扫描
func managedInspectScan(ctx *ScanContext) error {
	dir := ctx.ProjectDir
	startTime := time.Now()
	logger.Info.Println("Auto scan dir:", dir)
	for _, inspector := range managedInspector {
		filepath.WalkDir(ctx.ProjectDir, func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}
			if !d.IsDir() {
				return nil
			}
			{
				s, e := filepath.Rel(ctx.ProjectDir, path)
				if strings.Count(filepath.ToSlash(s), "/") > 5 || e != nil {
					return filepath.SkipDir
				}
			}
			if inspector.CheckDir(path) {
				logger.Debug.Println("Matched", inspector, path)
				rs, e := inspector.Inspect(path)
				if e != nil {
					logger.Info.Println("inspect failed.", inspector.String(), e.Error())
					logger.Debug.Printf("%+v\n", e)
					if e := base.UnwrapToInspectorError(e); e != nil {
						ctx.InspectorError = append(ctx.InspectorError, *e)
					}
				} else {
					for _, it := range rs {
						ctx.AddManagedModule(it)
					}
				}
				return filepath.SkipDir
			}
			return nil
		})
	}
	endTime := time.Now()
	logger.Info.Println("Scan terminated. Cost time:", endTime.Sub(startTime))
	logger.Info.Println("Total modules:", len(ctx.ManagedModules))
	return nil
}
