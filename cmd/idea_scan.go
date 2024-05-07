package cmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
				fmt.Println(model.GenerateIdeaErrorOutput(model.IdeaScanDirInvalid))
				SetGlobalExitCode(1)
				return
			}
			if e := initLogger(); e != nil {
				fmt.Println(model.GenerateIdeaErrorOutput(model.IdeaLogFileCreateFailed))
				SetGlobalExitCode(1)
				return
			}
			task := model.CreateScanTask(dir, model.TaskKindNormal, model.TaskTypeIdea)
			task.ProjectId = ProjectId
			ctx := model.WithScanTask(context.TODO(), task)
			if e := inspector.Scan(ctx); e != nil {
				fmt.Println(model.GenerateIdeaErrorOutput(e))
				SetGlobalExitCode(1)
				return
			}
			fmt.Println(model.GenerateIdeaOutput(ctx))
		},
	}
	c.Flags().StringVar(&dir, "dir", "", "project base dir")
	c.Flags().StringVar(&env.Scope, "scope", "", "")
	// compatible with recently version
	var indexOfIDEASCAN = slices.Index(os.Args, "ideascan")
	if indexOfIDEASCAN > -1 && len(os.Args) > indexOfIDEASCAN+1 && !strings.HasPrefix(os.Args[indexOfIDEASCAN+1], "-") {
		dir = os.Args[indexOfIDEASCAN+1]
		c.Args = cobra.ExactArgs(1)
	} else {
		c.Args = cobra.NoArgs
		must.Must(c.MarkFlagRequired("dir"))
		must.Must(c.MarkFlagDirname("dir"))
	}
	c.Flags().StringVar(&ProjectId, "project-id", "", "team id")
	c.Flags().StringVar(&env.GradleProjects, "gradle-projects", "", "")
	return c
}
