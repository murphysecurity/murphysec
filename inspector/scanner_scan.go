package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
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
	Logger.Info("File hash scanning")
	if e := FileHashScan(ctx); e != nil {
		Logger.Error("FileHash calc failed", zap.Error(e))
	}
	voModules := make([]api.VoModule, 0)
	for _, it := range task.Modules {
		voModules = append(voModules, api.VoModule{
			Dependencies:   it.Dependencies,
			FileHashList:   nil,
			Language:       it.Language,
			Name:           it.Name,
			PackageManager: it.PackageManager,
			RelativePath:   it.RelativePath,
			RuntimeInfo:    it.RuntimeInfo,
			Version:        it.Version,
			ModuleUUID:     it.UUID,
			ModuleType:     api.ModuleTypeVersion,
			ScanStrategy:   string(it.ScanStrategy),
		})
	}
	list := make([]api.VoFileHash, 0)
	for _, it := range task.FileHashes {
		for _, hash := range it.Hash {
			list = append(list, api.VoFileHash{
				Path: it.Path,
				Hash: hash,
			})
		}
	}
	voModules = append(voModules, api.VoModule{
		FileHashList: list,
		Language:     model.Cxx,
		ModuleType:   api.ModuleTypeFileHash,
		ModuleUUID:   _CPPModuleUUID,
	})
	fmt.Println(string(must.A(json.Marshal(voModules))))
}
