package scan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/common"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"path/filepath"
)

func Cmd() *cobra.Command {
	var c cobra.Command
	c.Use = "scan"
	c.Args = cobra.ExactArgs(1)
	c.Run = scanRun
	return &c
}

func scanRun(cmd *cobra.Command, args []string) {
	var (
		// workaround
		ctx     = ui.With(context.TODO(), ui.CLI{})
		scanDir = args[0]
		e       error
	)
	// get absolute path and check if a directory
	scanDir, e = filepath.Abs(scanDir)
	if e != nil {
		cv.DisplayScanInvalidPath(ctx, e)
	}
	if !utils.IsDir(scanDir) {
		cv.DisplayScanInvalidPathMustDir(ctx, nil)
		exitcode.Set(1)
		return
	}

	// init logging
	ctx, e = common.InitLogger(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		exitcode.Set(1)
		return
	}
	var logger = logctx.Use(ctx).Sugar()

	// init API
	e = common.InitAPIClient(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		logger.Error(e)
		exitcode.Set(1)
		return
	}

	_, e = scan(ctx, scanDir, model.AccessTypeCli)
	if e != nil {
		logger.Error(e)
		exitcode.Set(1)
		return
	}
}

func IdeaScan() *cobra.Command {
	var c cobra.Command
	c.Use = "ideascan"
	c.Args = cobra.ExactArgs(1)
	c.Run = ideascanRun
	c.Hidden = true
	c.Flags().String("ide", "", "unused")
	must.M(c.Flags().MarkHidden("ide"))
	return &c
}

func ideascanRun(cmd *cobra.Command, args []string) {
	var (
		// workaround
		ctx     = ui.With(context.TODO(), ui.None{})
		scanDir = args[0]
		e       error
	)
	// get absolute path and check if a directory
	scanDir, e = filepath.Abs(scanDir)
	if e != nil {
		cv.DisplayScanInvalidPath(ctx, e)
	}
	if !utils.IsDir(scanDir) {
		cv.DisplayScanInvalidPathMustDir(ctx, nil)
		exitcode.Set(1)
		return
	}

	// init logging
	ctx, e = common.InitLogger(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		exitcode.Set(1)
		return
	}
	var logger = logctx.Use(ctx).Sugar()

	// init API
	e = common.InitAPIClient(ctx)
	if e != nil {
		cv.DisplayInitializeFailed(ctx, e)
		logger.Error(e)
		exitcode.Set(1)
		return
	}

	task, e := scan(ctx, scanDir, model.AccessTypeIdea)
	if e != nil {
		logger.Error(e)
		exitcode.Set(1)
		return
	}
	fmt.Println(string(must.A(json.MarshalIndent(model.GetIDEAOutput(task), "", "  "))))
}
