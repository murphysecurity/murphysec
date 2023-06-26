package internalcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/common"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/inspector"
	logger2 "github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

func scannerScanCmd() *cobra.Command {
	var c cobra.Command
	c.Use = "scanner_scan <DIR>"
	c.Args = cobra.ExactArgs(1)
	c.Run = scannerScanRun

	return &c
}

func scannerScanRun(cmd *cobra.Command, args []string) {
	var (
		ctx     = ui.With(context.TODO(), ui.None)
		scanDir = args[0]
		e       error
	)
	env.ScannerScan = true
	common.LogLevel = logger2.LevelDebug
	ctx, e = common.InitLogger(ctx)
	if e != nil {
		fmt.Fprintf(os.Stderr, "init logger failed: %v\n", e)
	}
	var logger = logctx.Use(ctx).Sugar()

	// get absolute path and check if a directory
	scanDir = must.A(filepath.Abs(scanDir))
	if !utils.IsDir(scanDir) {
		logger.Error("not a directory")
		exitcode.Set(1)
		return
	}

	var scantask = &model.ScanTask{
		Ctx:         ctx,
		ProjectPath: scanDir,
		AccessType:  model.AccessTypeCli,
		Mode:        model.ScanModeSource,
		TaskId:      "",
		SubtaskId:   "",
		Modules:     nil,
		Result:      nil,
	}
	ctx = model.WithScanTask(ctx, scantask)
	e = inspector.ManagedInspect(ctx)
	if e != nil {
		logger.Error(e)
		exitcode.Set(1)
	}

	type wrapper struct {
		Modules                            []model.Module                `json:"modules"`
		ComponentCodeFragment              []model.ComponentCodeFragment `json:"component_code_fragment"`
		ScannerShouldEnableMavenBackupScan bool                          `json:"scanner_should_enable_maven_backup_scan"`
	}
	w := wrapper{
		Modules:                            utils.NoNilSlice(scantask.Modules),
		ComponentCodeFragment:              utils.NoNilSlice(scantask.CodeFragments),
		ScannerShouldEnableMavenBackupScan: env.ScannerShouldEnableMavenBackupScan,
	}
	time.Sleep(time.Second) // todo: workaround
	fmt.Println(string(must.M1(json.MarshalIndent(w, "", "  "))))
}
