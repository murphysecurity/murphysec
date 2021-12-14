package scanner

import (
	"fmt"
	"murphysec-cli-simple/api"
)

func PrintDetectReport(result *api.ScanResult) {
	fmt.Printf("扫描了 %d 个依赖项，找到了 %d 个问题\n", result.DependenciesCount, result.IssuesCount)
	ilc := result.IssuesLevelCount
	fmt.Printf("其中，关键项：%d 高危项：%d 中危项：%d 低微项：%d", ilc.Critical, ilc.High, ilc.Medium, ilc.Low)
}
