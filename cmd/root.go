package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/plugin"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/version"
	"os"
)

var (
	showVersion bool
)

func RootCmd() *cobra.Command {
	return rootCmd()
}

func rootCmd() *cobra.Command {
	selfName := "murphysec-cli"
	if len(os.Args) > 0 {
		selfName = os.Args[0]
	}

	c := &cobra.Command{
		Use:   selfName,
		Short: "murphysec-cli : An open source component security detection tool.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			for _, it := range plugin.Plugins {
				output.Debug(fmt.Sprintf("plugin: %v", it.Info().Name))
			}
			if showVersion {
				version.PrintVersionInfo()
				os.Exit(0)
			}
		},
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	c.PersistentFlags().BoolVarP(&showVersion, "version", "", false, "output version information and exit")
	c.PersistentFlags().BoolVarP(&output.Colorful, "color", "", true, "colorize the output")
	c.PersistentFlags().BoolVarP(&output.Verbose, "verbose", "v", false, "show verbose log")
	c.PersistentFlags().StringVarP(&conf.APITokenCliOverride, "token", "", "", "specify the API token")
	c.AddCommand(authCmd())
	c.AddCommand(scanCmd())
	return c
}

func rootHandler(cmd *cobra.Command, args []string) {
	must.Must(cmd.Help())
}
