package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/plugin"
	"murphysec-cli-simple/util/output"
	"os"
)

var scanDir string

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan",
		Run: func(cmd *cobra.Command, args []string) {
			for _, it := range plugin.Plugins {
				output.Info(fmt.Sprintf("Try match project by: %s", it.Info().Name))
				if it.MatchPath(scanDir) {
					output.Info(fmt.Sprintf("Match project succeed: %s", it.Info().Name))
					if e := scanByPlugin(it, scanDir); e != nil {
						output.Error(e.Error())
						os.Exit(-1)
					}
					return
				}
			}
			output.Error("Unable to inspect current directory, you can specify a directory by --dir <dir>")
			os.Exit(-1)
		},
		TraverseChildren: true,
	}
	c.PersistentFlags().StringVarP(&scanDir, "dir", "d", ".", "project root dir")
	for i := range plugin.Plugins {
		p := plugin.Plugins[i]
		pc := &cobra.Command{
			Use:              p.Info().Name,
			Short:            p.Info().ShortDescription,
			TraverseChildren: true,
			Run: func(cmd *cobra.Command, args []string) {
				if e := scanByPlugin(p, scanDir); e != nil {
					output.Error(e.Error())
					os.Exit(-1)
				}
			},
		}
		p.SetupScanCmd(pc)
		c.AddCommand(pc)
	}
	return c
}
