package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/env"
	"murphysec-cli-simple/inspector"
	"murphysec-cli-simple/utils"
)

var CliJsonOutput bool

var DeepScan bool
var ProjectId string

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan DIR",
		Run: scanRun,
	}
	c.Flags().BoolVar(&CliJsonOutput, "json", false, "json output")
	if env.AllowDeepScan {
		c.Flags().BoolVar(&DeepScan, "deep", false, "deep scan, will upload the source code")
	}
	c.Flags().StringVar(&ProjectId, "project-id", "", "team id")
	c.Args = cobra.ExactArgs(1)
	return c
}

var IotScan bool

func binscanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "binscan DIR",
		Run: binscanRun,
	}
	c.Flags().BoolVar(&IotScan, "iot", false, "enable iot scan")
	c.Args = cobra.ExactArgs(1)
	return c
}

func binscanRun(cmd *cobra.Command, args []string) {
	path := args[0]
	if !utils.IsPathExist(path) {
		fmt.Println("路径不存在")
		SetGlobalExitCode(1)
		return
	}
	t := api.TaskKindBinary
	if IotScan {
		t = api.TaskKindIotScan
	}
	ctx := inspector.NewBinaryScanContext(path, t)
	if e := inspector.BinScan(ctx); e != nil {
		SetGlobalExitCode(1)
	}
}
