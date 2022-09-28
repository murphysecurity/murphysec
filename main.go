package main

import (
	"fmt"
	"github.com/murphysecurity/murphysec/cmd"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/utils"
	"os"
	"time"
)

func main() {
	var (
		timeNow = time.Now()
	)
	e := cmd.Execute()
	if e != nil {
		os.Exit(-1)
	}
	display.CLI.Display(display.MsgNotice, fmt.Sprintf("任务耗时:%vs", utils.SubtractTime(timeNow, time.Now())))
	os.Exit(cmd.GetGlobalExitCode())
}
