package scanner

import (
	"fmt"
	"murphysec-cli-simple/api"
)

func PrintDetectReport(result *api.ScanResult) {
	fmt.Printf("扫描了 %d 个依赖项，找到了 %d 个问题\n", result.DependenciesCount, result.IssuesCount)
	ilc := result.IssuesLevelCount
	fmt.Println(fmt.Sprintf("其中，关键项：%d 高危项：%d 中危项：%d 低危项：%d", ilc.Critical, ilc.High, ilc.Medium, ilc.Low))
	fmt.Println(fmt.Sprintf("任务结果详情请点击：https://www.murphysec.com/control/code-scan/%s", result.TaskId))
}
