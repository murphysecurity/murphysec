package cmd

import (
	"github.com/spf13/cobra"
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

}
