package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/base"
	"murphysec-cli-simple/inspector"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
)

func ideaScanCmd() *cobra.Command {
	var dir string
	c := &cobra.Command{
		Hidden: true,
		Use:    "ideascan --dir ProjectDir",
		Run:    ideascanRun,
	}
	c.Flags().StringVar(&dir, "dir", "", "project base dir")
	c.Args = cobra.NoArgs
	must.Must(c.MarkFlagRequired("dir"))
	must.Must(c.MarkFlagDirname("dir"))
	c.Flags().StringVar(&ProjectId, "project-id", "", "team id")
	return c
}

func ideascanRun(cmd *cobra.Command, args []string) {
	dir := must.String(cmd.Flags().GetString("dir"))
	ctx, e := inspector.NewTaskContext(dir, base.TaskTypeIdea)
	if e != nil {
		logger.Err.Println(e)
		reportIdeaErr(IdeaScanDirInvalid, "")
		SetGlobalExitCode(1)
		return
	}
	ctx.ProjectId = ProjectId
	_, e = inspector.Scan(ctx)
	if e != nil {
		reportIdeaErr(e, "")
		SetGlobalExitCode(3)
		return
	}
	fmt.Println(string(must.Byte(json.MarshalIndent(generatePluginOutput(ctx), "", "  "))))
}
