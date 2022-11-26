package internalcmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"path/filepath"
)

func internalReadTaskIdCmd() *cobra.Command {
	var logger = must.A(zap.NewDevelopment())
	var ctx = logctx.With(context.TODO(), logger)

	var c cobra.Command
	c.Use = "read-task-id <DIR>"
	c.Flags().String("type", "", "")
	must.M(cobra.MarkFlagRequired(c.Flags(), "type"))
	c.Args = cobra.ExactArgs(1)

	c.Run = func(cmd *cobra.Command, args []string) {
		var acct = model.AccessType(cmd.Flag("type").Value.String())
		if c, e := config.ReadRepoConfig(ctx, must.A(filepath.Abs(args[0])), acct); e != nil {
			logger.Error(e.Error())
			exitcode.Set(1)
			return
		} else {
			fmt.Println(c.TaskId)
		}
	}

	return &c
}
