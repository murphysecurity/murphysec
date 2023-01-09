package cmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/envinspection"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/spf13/cobra"
)

func envScanCmd() *cobra.Command {
	var projName string
	c := &cobra.Command{
		Use: "envscan",
		Run: func(cmd *cobra.Command, args []string) {
			initConsoleLoggerOrExit()
			if e := envinspection.InspectEnv(utils.WithLogger(context.TODO(), LOG), projName); e != nil {
				fmt.Println(e.Error())
				SetGlobalExitCode(1)
			}
		},
	}

	c.Flags().StringVar(&projName, "project-name", "", "")
	c.Flags().StringVar(&utils.NetworkInterfaceName, "interface", "", "Only used in default project identifier")

	return c
}
