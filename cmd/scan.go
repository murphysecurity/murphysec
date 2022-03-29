package cmd

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/env"
)

var CliJsonOutput bool

var DeepScan bool

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan DIR",
		Run: scanRun,
	}
	c.Flags().BoolVar(&CliJsonOutput, "json", false, "json output")
	if env.AllowDeepScan {
		c.Flags().BoolVar(&DeepScan, "deep", false, "deep scan, will upload the source code")
	}
	c.Args = cobra.ExactArgs(1)
	return c
}
