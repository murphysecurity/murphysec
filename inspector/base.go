package inspector

import (
	"murphysec-cli-simple/env"
)

func shouldUploadFile(ctx *ScanContext) bool {
	if !env.AllowDeepScan {
		return false
	}
	if len(ctx.ManagedModules) == 0 {
		return true
	}
	for _, it := range ctx.ManagedModules {
		if it.PackageManager == "maven" {
			return true
		}
	}
	return false
}

//ctx, e := CreateTaskContext(must.String(filepath.Abs(dir)), api.TaskTypeIdea)
//if e != nil {
//	fmt.Println(e)
//	panic(e)
//}
//if e := managedInspectScan(ctx); e != nil {
//	logger.Err.Println("Managed inspect failed.", e.Error())
//	logger.Debug.Printf("%+v", e)
//}
//if ctx.ManagedModules == nil {
//	ctx.ManagedModules = []base.Module{}
//}
//fmt.Println(string(must.Byte(json.Marshal(ctx.ManagedModules))))

func ScannerScan(dir string) {
}
