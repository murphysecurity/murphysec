package inspector

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"time"
)

func ManagedInspect(ctx context.Context) error {
	var logger = logctx.Use(ctx)
	scanTask := model.UseScanTask(ctx)

	// 获取扫描路径
	baseDir := scanTask.ProjectPath
	logger.Info("Auto scan dir", zap.String("dir", baseDir))

	// 扫
	var scanner = &dirScanner{
		inspectors: module.Inspectors,
		root:       baseDir,
	}
	scanner.scan()

	logger.Sugar().Infof("Found %d directories", len(scanner.scannedDirs))

	// 对扫到的内容，逐个开始检查
	for idx, it := range scanner.scannedDirs {
		st := time.Now()
		// 创建检查任务
		c := model.WithInspectionTask(ctx, scanTask.BuildInspectionTask(it.path))
		// 绑定 logger
		c = logctx.With(c, logger.Named(fmt.Sprintf("%s-%d", it.inspector.String(), idx)))

		// Do!
		logger.Sugar().Infof("Begin: %s", it.String())
		e := it.inspector.InspectProject(c)
		logger.Sugar().Infof("End: %s, duration: %v", it.String(), time.Since(st))
		if e != nil {
			logger.Error("InspectError", zap.Error(e), zap.Any("inspector", it))
		}
	}

	var components []model.Component
	for _, m := range scanTask.Modules {
		components = append(components, m.ComponentList()...)
	}
	logger.Sugar().Debugf("scan code fragments, total %d components", len(components))
	components = utils.DistinctSlice(components)
	previews, e := scanFragment(ctx, scanTask.ProjectPath, components)
	if e != nil {
		logger.Sugar().Errorf("errors during scan fragment: %s", e.Error())
	} else {
		scanTask.CodeFragments = previews
	}
	return nil
}
