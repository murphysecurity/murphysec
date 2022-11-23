package cmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/build_flags"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/module"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var versionFlag bool
var cliServerAddressOverride string
var cliTokenOverride string
var allowInsecure bool
var logLevel logger.Level

var rootCtx = context.TODO()

func preRun(cmd *cobra.Command, args []string) error {
	if versionFlag {
		fmt.Println(version.FullInfo())
		fmt.Printf("Supported modules: %s\n", strings.Join(module.GetSupportedModuleList(), ", "))
		os.Exit(0)
	}

	logger.LogFileCleanup()

	if e := logLevel.Of(consoleLogLevelOverride); e != nil {
		return e
	}

	return nil
}

func rootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:               "murphysec",
		PersistentPreRunE: preRun,
		TraverseChildren:  true,
		Run: func(cmd *cobra.Command, args []string) {
			must.Must(cmd.Help())
		},
	}

	// Version
	c.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "show version and exit")

	// Logging
	c.PersistentFlags().BoolVar(&disableLogFile, "no-log-file", false, "do not write log file")
	c.PersistentFlags().StringVar(&cliLogFilePathOverride, "write-log-to", "", "specify log file path")
	c.PersistentFlags().StringVar(&consoleLogLevelOverride, "log-level", "silent", "specify log level, must be silent|error|warn|info|debug")
	c.PersistentFlags().BoolVar(&enableNetworkLog, "network-log", false, "print network data")

	// API: Authentication & Network
	c.PersistentFlags().StringVar(&cliTokenOverride, "token", "", "specify API token")
	c.PersistentFlags().StringVar(&cliServerAddressOverride, "server", "", "specify server address")
	c.PersistentFlags().BoolVarP(&allowInsecure, "allow-insecure", "x", false, "Allow insecure TLS connection")

	c.AddCommand(authCmd())
	c.AddCommand(scanCmd())
	if build_flags.InternalFeature {
		c.AddCommand(internalCmd())
	}
	c.AddCommand(machineCmd())
	//c.AddCommand(envScanCmd())
	return c
}

func Execute() error {
	return rootCmd().Execute()
}
