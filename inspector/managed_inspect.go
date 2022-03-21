package inspector

import (
	"io/fs"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/go_mod"
	"murphysec-cli-simple/module/gradle"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/module/npm"
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
}

// 受管理扫描
// returns ErrNoEngineMatched ErrNoModule
func managedInspectScan(ctx *ScanContext) error {
	dir := ctx.ProjectDir
	startTime := time.Now()
	logger.Info.Println("Auto scan dir:", dir)
	for _, inspector := range managedInspector {
		logger.Debug.Println("For:", inspector.String())
		filepath.WalkDir(ctx.ProjectDir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return nil
			}
			{
				s, e := filepath.Rel(ctx.ProjectDir, path)
				if strings.Count(filepath.ToSlash(s), "/") > 3 || e != nil {
					return filepath.SkipDir
				}
			}
			logger.Debug.Println("Visit dir:", path)
			if inspector.CheckDir(path) {
				logger.Debug.Println("Matched")
				rs, e := inspector.Inspect(path)
				if e != nil {
					logger.Info.Println("inspect failed.", e.Error())
					logger.Debug.Printf("%+v\n", e)
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
	if len(ctx.ManagedModules) < 1 {
		return ErrNoModule
	}
	return nil
}
