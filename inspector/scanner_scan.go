package inspector

import (
	"encoding/json"
	"fmt"
	base2 "murphysec-cli-simple/base"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
)

func ScannerScan(dir string) {
	ctx, e := NewTaskContext(must.String(filepath.Abs(dir)), base2.TaskTypeCli)
	if e != nil {
		panic(e)
	}
	if e := managedInspectScan(ctx); e != nil {
		logger.Err.Println("Managed inspect failed.", e.Error())
		logger.Debug.Printf("%+v", e)
	}
	if ctx.ManagedModules == nil {
		ctx.ManagedModules = []base.Module{}
	}
	fmt.Println(string(must.Byte(json.Marshal(ctx.ManagedModules))))
}
