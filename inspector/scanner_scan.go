package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
)

func ScannerScan(dir string) {
	task := model.CreateScanTask(must.A(filepath.Abs(dir)), model.TaskKindNormal, model.TaskTypeIdea)
	ctx := model.WithScanTask(context.TODO(), task)
	if e := managedInspect(ctx); e != nil {
		logger.Err.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%+v", e)
	}
	if task.Modules == nil {
		task.Modules = []model.Module{}
	}
	fmt.Println(string(must.A(json.Marshal(task.Modules))))
}
