package cmd

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/inspector"
)

func scannerCmd() *cobra.Command {
	c := &cobra.Command{
		Use:    "scanner_scan",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			dir := args[0]
			inspector.ScannerScan(dir)
		},
	}
	c.Args = cobra.ExactArgs(1)
	return c
}
