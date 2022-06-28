package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"go.uber.org/zap"
	"path/filepath"
)

func ScannerScan(dir string) {
	task := model.CreateScanTask(must.A(filepath.Abs(dir)), model.TaskKindNormal, model.TaskTypeIdea)
	ctx := model.WithScanTask(context.TODO(), task)
	if e := managedInspect(ctx); e != nil {
		Logger.Error("Managed inspect error", zap.Error(e))
	}
	if task.Modules == nil {
		task.Modules = []model.Module{}
	}
	fmt.Println(string(must.A(json.Marshal(task.Modules))))
}
