package cmd

import (
	"context"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/spf13/cobra"
)

func dockerScanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "dockerfile",
		Run: func(cmd *cobra.Command, args []string) {
			initConsoleLoggerOrExit()
			ctx := context.TODO()
			task := model.CreateScanTask(args[0], model.TaskKindDockerfile, model.TaskTypeCli)
			ctx = model.WithScanTask(ctx, task)
			inspector.InspectDockerfile(ctx)
		},
	}
	c.Args = cobra.ExactArgs(1)
	return c
}
