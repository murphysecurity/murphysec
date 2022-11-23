package cmd

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/spf13/cobra"
)

func scanCmd() *cobra.Command {
	var c cobra.Command
	c.Use = "scan"
	c.Args = cobra.ExactArgs(1)
	c.Run = scanRun
	return &c
}

func scanRun(cmd *cobra.Command, args []string) {
	var ctx = commonInitCLI()
	var scanDir = args[0]
	if !utils.IsDir(scanDir) {
		displayScanInvalidDir(ctx)
		exitcode.Set(1)
		return
	}
}

func commonInitCLI() context.Context {
	mustInitLogger()
	var ctx = rootCtx
	ctx = ui.With(rootCtx, ui.CLI{})
	e := initAPI(ctx)
	if e != nil {
		exitcode.Set(1)
		return nil
	}
	return ctx
}
