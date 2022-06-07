package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"path/filepath"
)

func ideaScanCmd() *cobra.Command {
	var dir string
	c := &cobra.Command{
		Hidden: true,
		Use:    "ideascan --dir ProjectDir",
		Run: func(cmd *cobra.Command, args []string) {
			if !filepath.IsAbs(dir) {
				dir = must.A(filepath.Abs(dir))
			}
			if !utils.IsPathExist(dir) {
				reportIdeaErr(IdeaScanDirInvalid, "")
				SetGlobalExitCode(1)
				return
			}
			logger.InitLogger()
			task := model.CreateScanTask(dir, model.TaskKindNormal, model.TaskTypeIdea)
			task.ProjectId = ProjectId
			ctx := model.WithScanTask(context.TODO(), task)
			if e := inspector.Scan(ctx); e != nil {
				reportIdeaErr(e, "")
				SetGlobalExitCode(3)
				return
			}
			fmt.Println(string(must.A(json.MarshalIndent(generatePluginOutput(ctx), "", "  "))))
		},
	}
	c.Flags().StringVar(&dir, "dir", "", "project base dir")
	c.Args = cobra.NoArgs
	must.Must(c.MarkFlagRequired("dir"))
	must.Must(c.MarkFlagDirname("dir"))
	c.Flags().StringVar(&ProjectId, "project-id", "", "team id")
	return c
}
