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
		Use:   "scan DIR",
		Run:   scanRun,
		Short: "Scan the source code of the specified project, currently supporting java, javascript, go, and python",
	}
	c.Flags().BoolVar(&CliJsonOutput, "json", false, "json output")
	if env.AllowDeepScan {
		c.Flags().BoolVar(&DeepScan, "deep", false, "deep scan, will upload the source code")
	}
	c.Flags().StringVar(&ProjectId, "project-id", "", "team id")
	c.Args = cobra.ExactArgs(1)
	return c
}

func binScanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "binscan DIR",
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			if !utils.IsPathExist(path) {
				fmt.Println("路径不存在")
				SetGlobalExitCode(1)
				return
			}
			ctx := inspector.NewBinaryScanContext(path, api.TaskKindBinary)
			if e := inspector.BinScan(ctx); e != nil {
				SetGlobalExitCode(1)
			}
		},
		Short: "Scan specified binary files and software artifacts, currently supporting .jar, .war, and common binary file formats (The file will be uploaded to the server for analysis.)",
	}
	c.Args = cobra.ExactArgs(1)
	return c
}

func iotScanCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "iotscan DIR",
		Short: "Scan the specified IoT device firmware, currently supporting .bin or other formats (The file will be uploaded to the server for analysis.)",
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			if !utils.IsPathExist(path) {
				fmt.Println("路径不存在")
				SetGlobalExitCode(1)
				return
			}
			ctx := inspector.NewBinaryScanContext(path, api.TaskKindIotScan)
			if e := inspector.BinScan(ctx); e != nil {
				SetGlobalExitCode(1)
			}
		},
	}
	c.Args = cobra.ExactArgs(1)
	return c
}
