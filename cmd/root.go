package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/version"
	"os"
)

var versionFlag bool
var managedMode bool
var logRedirectPath string
var noLogFile bool
var taskInfo string

func rootCmd() *cobra.Command {
	argsMap := map[string]bool{}
	for _, it := range os.Args {
		argsMap[it] = true
	}
	c := &cobra.Command{
		Use:               "murphysec",
		PersistentPreRunE: preRun,
		TraverseChildren:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	c.PersistentFlags().BoolVar(&versionFlag, "version", false, "show version and exit")
	c.PersistentFlags().BoolVar(&noLogFile, "no-log-file", false, "do not write log file")
	c.PersistentFlags().StringVar(&logRedirectPath, "write-log-to", "", "specify log file path")
	// workaround avoid err
	if argsMap["--managed-mode"] {
		c.PersistentFlags().BoolVar(&managedMode, "managed-mode", false, "")
	}
	if argsMap["--task-info"] {
		c.PersistentFlags().StringVar(&taskInfo, "task-info", "", "")
	}
	c.AddCommand(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello world!")
		},
	})
	c.AddCommand(authCmd())
	return c
}

func preRun(cmd *cobra.Command, args []string) error {
	if versionFlag {
		version.PrintVersionInfo()
		return nil
	}
	if managedMode {
		logger.InitManagedMode()
	}
	if !noLogFile {
		must.Must(logger.InitFileLog(logRedirectPath))
	}
	return nil
}

func Execute() error {
	return rootCmd().Execute()
}
