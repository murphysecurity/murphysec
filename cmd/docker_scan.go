package cmd

import (
	"context"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"path/filepath"
)

func dockerScanCmd() *cobra.Command {
	var jsonOutput bool
	c := &cobra.Command{
		Use: "dockerfile",
		Run: func(cmd *cobra.Command, args []string) {
			initConsoleLoggerOrExit()
			ctx := context.TODO()
			taskType := model.TaskTypeCli
			if jsonOutput {
				taskType = model.TaskTypeJenkins
			}
			task := model.CreateScanTask(must.A(filepath.Abs(args[0])), model.TaskKindDockerfile, taskType)
			ctx = model.WithScanTask(ctx, task)
			_ = inspector.InspectDockerfile(ctx)
		},
	}
	c.Flags().BoolVar(&jsonOutput, "json", false, "json output")
	c.Args = cobra.ExactArgs(1)
	return c
}
