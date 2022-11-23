package cmd

import "github.com/spf13/cobra"

func internalCmd() *cobra.Command {
	var c = &cobra.Command{Use: "internal", Hidden: true}
	c.AddCommand(cppFileHashCmd())
	c.AddCommand(internalWriteTaskId())
	c.AddCommand(internalReadTaskId())
	return c
}
