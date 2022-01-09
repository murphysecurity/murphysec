package cmd

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/inspector"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils/must"
)

func ideaScanCmd() *cobra.Command {
	var dir string
	var packageManager string
	c := &cobra.Command{
		Use: "ideascan --dir ProjectDir --package-manager PackageManagerType",
		Run: func(cmd *cobra.Command, args []string) {
			var pmType = base.PackageManagerTypeOfName(packageManager)
			_, e := inspector.IdeaScan(dir, pmType)
			if e != nil {
				SetGlobalExitCode(1)
				logger.Err.Println("idea plugin scan failed.", e.Error())
			}
		},
	}
	c.Flags().StringVar(&dir, "dir", "", "project base dir")
	c.Flags().StringVar(&packageManager, "package-manager", "", "package manager type, must be maven|gomod")
	c.Args = cobra.NoArgs
	must.Must(c.MarkFlagRequired("dir"))
	must.Must(c.MarkFlagDirname("dir"))
	must.Must(c.MarkFlagRequired("package-manager"))
	return c
}
