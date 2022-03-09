package cmd

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/inspector"
	"murphysec-cli-simple/utils/must"
	"path/filepath"
)

func binScan() *cobra.Command {
	c := &cobra.Command{
		Use:  "binscan",
		Args: cobra.ExactArgs(1),
		Run:  binScanRun,
	}
	return c
}

func binScanRun(cmd *cobra.Command, args []string) {
	fp := args[0]
	inspector.BinScan(must.String(filepath.Abs(fp)))
}
