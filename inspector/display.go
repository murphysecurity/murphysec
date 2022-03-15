package inspector

import (
	"fmt"
	"murphysec-cli-simple/api"
)

func displayTaskCreating(ctx *ScanContext) {
	if ctx.TaskType == api.TaskTypeCli {
		fmt.Println("正在创建扫描任务，请稍候，项目名称：", ctx.ProjectName)
	}
}

func displayTaskCreated(ctx *ScanContext) {
	if ctx.TaskType == api.TaskTypeCli {
		fmt.Println("扫描任务已创建")
	}
}

func displayManagedScanning(ctx *ScanContext) {
	if ctx.TaskType == api.TaskTypeCli {
		fmt.Println("正在执行扫描")
	}
}
