package cmd

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
	"os"
)

func internalWriteTaskId() *cobra.Command {
	var logger = must.A(zap.NewDevelopment())
	var ctx = logctx.With(context.TODO(), logger)

	var c cobra.Command
	c.Use = "write-task-id"
	c.Flags().String("type", "", "")
	must.M(cobra.MarkFlagRequired(c.Flags(), "type"))
	c.Args = cobra.ExactArgs(1)

	c.Run = func(cmd *cobra.Command, args []string) {
		var acct = model.AccessType(cmd.Flag("type").Value.String())
		if e := config.WriteRepoConfig(ctx, mustGetCWD(), acct, config.RepoConfig{TaskId: args[0]}); e != nil {
			logger.Error(e.Error())
			exitcode.Set(1)
			return
		}
	}

	return &c
}

func internalReadTaskId() *cobra.Command {
	var logger = must.A(zap.NewDevelopment())
	var ctx = logctx.With(context.TODO(), logger)

	var c cobra.Command
	c.Use = "read-task-id"
	c.Flags().String("type", "", "")
	must.M(cobra.MarkFlagRequired(c.Flags(), "type"))
	c.Args = cobra.ExactArgs(1)

	c.Run = func(cmd *cobra.Command, args []string) {
		var acct = model.AccessType(cmd.Flag("type").Value.String())
		if c, e := config.ReadRepoConfig(ctx, mustGetCWD(), acct); e != nil {
			logger.Error(e.Error())
			exitcode.Set(1)
			return
		} else {
			fmt.Println(c.TaskId)
		}
	}

	return &c
}

func mustGetCWD() string {
	return must.A(os.Getwd())
}
