package internalcmd

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	var c = &cobra.Command{Use: "internal", Hidden: true}
	c.AddCommand(machineidCmd())
	c.AddCommand(scannerScanCmd())
	return c
}
