package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
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
	return c
}

func ideascanRun(cmd *cobra.Command, args []string) {
	dir := must.String(cmd.Flags().GetString("dir"))
	ctx, e := inspector.NewTaskContext(dir, base.TaskTypeIdea)
	if e != nil {
		logger.Err.Println(e)
		reportIdeaErr(IdeaScanDirInvalid)
		SetGlobalExitCode(1)
		return
	}
	_, e = inspector.Scan(ctx)
	if e != nil {
		logger.Err.Println(e)
		logger.Debug.Printf("%v+\n", e)
		if errors.Is(e, api.ErrTokenInvalid) {
			reportIdeaErr(IdeaTokenInvalid)
		} else if errors.Is(e, api.ErrServerRequest) {
			reportIdeaErr(IdeaServerRequestFailed)
		} else if errors.Is(e, api.ErrTimeout) {
			reportIdeaErr(IdeaApiTimeout)
		} else {
			reportIdeaErr(IdeaUnknownProject)
		}
		SetGlobalExitCode(3)
		return
	}
	fmt.Println(string(must.Byte(json.MarshalIndent(generatePluginOutput(ctx), "", "  "))))
}
