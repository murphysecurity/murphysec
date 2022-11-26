package internalcmd

import (
	"context"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"path/filepath"
)

func internalWriteTaskIdCmd() *cobra.Command {
	var logger = must.A(zap.NewDevelopment())
	var ctx = logctx.With(context.TODO(), logger)

	var c cobra.Command
	c.Use = "write-task-id <DIR> <TaskID>"
	c.Flags().String("type", "", "")
	must.M(cobra.MarkFlagRequired(c.Flags(), "type"))
	c.Args = cobra.ExactArgs(2)

	c.Run = func(cmd *cobra.Command, args []string) {
		var acct = model.AccessType(cmd.Flag("type").Value.String())
		if e := config.WriteRepoConfig(ctx, must.A(filepath.Abs(args[0])), acct, config.RepoConfig{TaskId: args[1]}); e != nil {
			logger.Error(e.Error())
			exitcode.Set(1)
			return
		}
	}

	return &c
}
