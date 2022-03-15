package inspector

import (
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/module/go_mod"
	"murphysec-cli-simple/module/gradle"
	"murphysec-cli-simple/module/maven"
	"murphysec-cli-simple/module/npm"
	"time"
)

var managedInspector = []base.Inspector{
	go_mod.New(),
	maven.New(),
	npm.New(),
	gradle.New(),
}

// 受管理扫描
// returns ErrNoEngineMatched ErrNoModule
func managedInspectScan(ctx *ScanContext) error {
	dir := ctx.ProjectDir
	startTime := time.Now()
	logger.Info.Println("Auto scan dir:", dir)
	var inspectors []base.Inspector
	{
		// 尝试匹配检测器
		logger.Debug.Println("Try match managed inspector...")
		for _, it := range managedInspector {
			if it.CheckDir(dir) {
				inspectors = append(inspectors, it)
			}
		}
		logger.Debug.Println("Matched managed inspector:", inspectors)
	}
	if len(inspectors) == 0 {
		logger.Debug.Println("No managed inspector matched")
		return ErrNoEngineMatched
	}

	for _, it := range inspectors {
		rs, e := it.Inspect(dir)
		if e != nil {
			logger.Err.Printf("Engine: %v scan failed. Reason: %+v\n", it, e)
			continue
		}
		logger.Info.Printf("Inspector terminated %v, total module: %v\n", it, len(rs))
		for _, it := range rs {
			ctx.AddManagedModule(it)
		}
	}
	endTime := time.Now()
	logger.Info.Println("Scan terminated. Cost time:", endTime.Sub(startTime))
	if len(ctx.ManagedModules) < 1 {
		return ErrNoModule
	}
	return nil
}
