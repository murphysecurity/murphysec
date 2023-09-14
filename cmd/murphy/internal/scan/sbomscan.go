package scan

import (
	"context"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/cv"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/spf13/cobra"
	"os"
)

func SbomScan() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "sbom <DIR>",
		Short: "Scan project dependencies and generate SBOM in SPDX-JSON",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var ctx = context.TODO()
			ctx = ui.With(ctx, ui.None)
			scanDir := args[0]
			scanDir, e := commonScanPreCheck(ctx, scanDir)
			if e != nil {
				return
			}
			ctx, e = commonInitNoAPI(ctx)
			if e != nil {
				return
			}
			logger := logctx.Use(ctx)
			spdx, e := processDir(ctx, scanDir)
			if e != nil {
				exitcode.Set(1)
				logger.Sugar().Error(e)
				return
			}
			if out == "" {
				out = "spdx.json"
			}
			e = os.WriteFile(out, spdx, 0644)
			if e != nil {
				exitcode.Set(1)
				logger.Sugar().Error(e)
				return
			}
		},
	}
	cmd.Flags().StringVar(&out, "out", "", "output file path")
	cmd.Flags().String("type", "", "")
	_ = cmd.Flags().MarkHidden("type")
	return cmd
}

func processDir(ctx context.Context, dir string) ([]byte, error) {
	var e error
	var task = model.ScanTask{
		Ctx:         ctx,
		ProjectPath: dir,
	}
	ctx = model.WithScanTask(ctx, &task)
	e = inspector.ManagedInspect(ctx)
	if e != nil {
		cv.DisplayScanFailed(ctx, e)
		return nil, e
	}
	spdxData := inspector.BuildSpdx(&task)
	return spdxData, nil
}
