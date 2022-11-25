package internalcmd

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	var c = &cobra.Command{Use: "internal", Hidden: true}
	c.AddCommand(cppFileHashCmd())
	c.AddCommand(internalWriteTaskIdCmd())
	c.AddCommand(internalReadTaskIdCmd())
	c.AddCommand(machineidCmd())
	return c
}
