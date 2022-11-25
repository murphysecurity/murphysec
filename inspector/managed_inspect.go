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
		c = utils.WithLogger(c, logger.Named(fmt.Sprintf("%s-%d", it.inspector.String(), idx)))

		// Do!
		logger.Sugar().Infof("Begin: %s, duration: %v", it.String(), time.Now().Sub(st))
		e := it.inspector.InspectProject(c)
		logger.Sugar().Infof("End: %s, duration: %v", it.String(), time.Now().Sub(st))
		if e != nil {
			logger.Error("InspectError", zap.Error(e), zap.Any("inspector", it))
		}
	}
	return nil
}
