package cmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"path/filepath"
)

var CliJsonOutput bool

var DeepScan bool
var ProjectId string

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan DIR",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()
			logger.InitLogger()
			projectDir := args[0]
			var e error
			if !filepath.IsAbs(projectDir) {
				projectDir, e = filepath.Abs(projectDir)
				if e != nil || !utils.IsPathExist(projectDir) {
					fmt.Println("读取路径失败", e.Error())
					SetGlobalExitCode(1)
					return
				}
			}
			tt := model.TaskTypeCli
			if CliJsonOutput {
				tt = model.TaskTypeJenkins
			}
			task := model.CreateScanTask(projectDir, model.TaskKindNormal, tt)
			task.ProjectId = ProjectId
			if env.SpecificProjectName != "" {
				task.ProjectName = env.SpecificProjectName
			}
			task.EnableDeepScan = DeepScan
			ctx = model.WithScanTask(ctx, task)

			if e := inspector.Scan(ctx); e != nil {
				if tt == model.TaskTypeJenkins {
					fmt.Println(model.GenerateIdeaErrorOutput(e))
				}
				SetGlobalExitCode(-1)
			} else {
				if tt == model.TaskTypeJenkins {
					fmt.Println(model.GenerateIdeaOutput(ctx))
				}
			}
		},
		Short: "Scan the source code of the specified project, currently supporting java, javascript, go, and python",
	}
	c.Flags().BoolVar(&CliJsonOutput, "json", false, "json output")
	if env.AllowDeepScan {
		c.Flags().BoolVar(&DeepScan, "deep", false, "deep scan, will upload the source code")
	}
	c.Flags().StringVar(&ProjectId, "project-id", "", "team id")
	must.Must(c.Flags().MarkHidden("project-id"))
	c.Flags().StringVar(&env.SpecificProjectName, "project-name", "", "force specific project name")
	c.Flags().BoolVar(&env.DisableGit, "skip-git", false, "force ignore git info")
	c.Args = cobra.ExactArgs(1)
	return c
}

func binScanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "binscan DIR",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()
			logger.InitLogger()
			projectDir := args[0]
			var e error
			if !filepath.IsAbs(projectDir) {
				projectDir, e = filepath.Abs(projectDir)
				if e != nil || !utils.IsPathExist(projectDir) {
					fmt.Println("读取路径失败", e.Error())
					SetGlobalExitCode(1)
					return
				}
			}
			task := model.CreateScanTask(projectDir, model.TaskKindBinary, model.TaskTypeCli)
			ctx = model.WithScanTask(ctx, task)
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
			ctx := context.TODO()
			logger.InitLogger()
			projectDir := args[0]
			var e error
			if !filepath.IsAbs(projectDir) {
				projectDir, e = filepath.Abs(projectDir)
				if e != nil || !utils.IsPathExist(projectDir) {
					fmt.Println("读取路径失败", e.Error())
					SetGlobalExitCode(1)
					return
				}
			}
			task := model.CreateScanTask(projectDir, model.TaskKindIotScan, model.TaskTypeCli)
			ctx = model.WithScanTask(ctx, task)
			if e := inspector.BinScan(ctx); e != nil {
				SetGlobalExitCode(1)
			}
		},
	}
	c.Args = cobra.ExactArgs(1)
	return c
}
