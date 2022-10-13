package cmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/envinspection"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/spf13/cobra"
)

func envScanCmd() *cobra.Command {
	return &cobra.Command{
		Use: "envscan",
		Run: func(cmd *cobra.Command, args []string) {
			initConsoleLoggerOrExit()
			if e := envinspection.InspectEnv(utils.WithLogger(context.TODO(), LOG)); e != nil {
				fmt.Println(e.Error())
				SetGlobalExitCode(1)
			}
		},
	}
}
