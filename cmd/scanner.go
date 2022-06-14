package cmd

import (
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
)

func scannerCmd() *cobra.Command {
	c := &cobra.Command{
		Use:    "scanner_scan",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			dir := args[0]
			must.Must(logger.InitLogger())
			inspector.ScannerScan(dir)
		},
	}
	c.Args = cobra.ExactArgs(1)
	return c
}
