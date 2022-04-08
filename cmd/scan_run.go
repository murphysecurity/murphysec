package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/base"
	"murphysec-cli-simple/display"
	"murphysec-cli-simple/inspector"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"strconv"
)

func scanRun(cmd *cobra.Command, args []string) {
	if CliJsonOutput {
		ctx, e := inspector.NewTaskContext(args[0], base.TaskTypeIdea)
		if e != nil {
			logger.Err.Println(e)
			reportIdeaErr(IdeaScanDirInvalid, "")
			SetGlobalExitCode(1)
			return
		}
		ctx.ProjectId = ProjectId
		_, e = inspector.Scan(ctx)
		if e != nil {
			reportIdeaErr(e, "")
			SetGlobalExitCode(3)
			return
		}
		fmt.Println(string(must.Byte(json.MarshalIndent(generatePluginOutput(ctx), "", "  "))))
		return
	}
	dir := args[0]
	source := base.TaskTypeCli
	if CliJsonOutput {
		source = base.TaskTypeJenkins
	} else {
		display.EnableANSI()
	}
	ui := source.UI()
	ctx, e := inspector.NewTaskContext(dir, source)
	if e != nil {
		logger.Err.Println(e)
		ui.Display(display.MsgError, "项目目录无效或不存在")
		SetGlobalExitCode(1)
		return
	}
	ctx.ProjectId = ProjectId
	ctx.EnableDeepScan = DeepScan
	if _, e = inspector.Scan(ctx); e != nil {
		ui.Display(display.MsgError, "扫描失败："+e.Error())
		SetGlobalExitCode(2)
		return
	}
	{
		totalDep := strconv.Itoa(ctx.ScanResult.DependenciesCount)
		totalVuln := strconv.Itoa(ctx.ScanResult.IssuesCompsCount)
		t := fmt.Sprint(
			"项目扫描成功，依赖数：",
			termenv.String(totalDep).Foreground(termenv.ANSIBrightCyan),
			"，漏洞数：",
			termenv.String(totalVuln).Foreground(termenv.ANSIBrightRed),
		)
		ui.Display(display.MsgNotice, t)
	}
}
