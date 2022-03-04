package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/version"
)

func machineCmd() *cobra.Command {
	return &cobra.Command{
		Use: "machine-id",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.MachineId())
		},
		Hidden: true,
	}
}
