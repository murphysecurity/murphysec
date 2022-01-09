package cmd

import (
	"github.com/spf13/cobra"
)

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan DIR",
		Run: scanRun,
	}
	c.Args = cobra.ExactArgs(1)
	return c
}

func scanRun(cmd *cobra.Command, args []string) {

}
