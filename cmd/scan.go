package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/inspector"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
)

var CliJsonOutput bool

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan DIR",
		Run: scanRun,
	}
	c.Flags().BoolVar(&CliJsonOutput, "json", false, "json output")
	c.Args = cobra.ExactArgs(1)
	return c
}

func scanRun(cmd *cobra.Command, args []string) {
	logger.Info.Println("CLI scan dir:", args[0], must.String(filepath.Abs(args[0])))
	var e error
	if CliJsonOutput {
		_, e = inspector.Scan(must.String(filepath.Abs(args[0])), api.TaskTypeJenkins)
	} else {
		_, e = inspector.Scan(must.String(filepath.Abs(args[0])), api.TaskTypeCli)
	}
	if e != nil {
		SetGlobalExitCode(1)
		if !CliJsonOutput {
			fmt.Printf("命令行扫描失败，错误：%v\n", e)
		}
		logger.Err.Printf("Cli scan failed. %v\n", e)
		logger.Debug.Printf("Cli scan failed. %+v\n", e)
	}
}
