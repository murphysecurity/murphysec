package cmd

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/inspector"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
)

var CliJsonOutput bool

var DeepScan bool

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan DIR",
		Run: scanRun,
	}
	c.Flags().BoolVar(&CliJsonOutput, "json", false, "json output")
	c.Flags().BoolVar(&DeepScan, "deep", false, "deep scan, will upload the source code")
	c.Args = cobra.ExactArgs(1)
	return c
}

func scanRun(cmd *cobra.Command, args []string) {
	logger.Info.Println("CLI scan dir:", args[0], must.String(filepath.Abs(args[0])))
	var e error
	println(filepath.Abs(args[0]))
	if CliJsonOutput {
		_, e = inspector.Scan(must.String(filepath.Abs(args[0])), api.TaskTypeJenkins, DeepScan)
	} else {
		_, e = inspector.Scan(must.String(filepath.Abs(args[0])), api.TaskTypeCli, DeepScan)
	}
	if e != nil {
		SetGlobalExitCode(1)
		logger.Err.Printf("Cli scan failed. %v\n", e)
		logger.Debug.Printf("Cli scan failed. %+v\n", e)
	}
}
