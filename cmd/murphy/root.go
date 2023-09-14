package murphy

import (
	"fmt"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/auth"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/binscan"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/common"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/internalcmd"
	"github.com/murphysecurity/murphysec/cmd/murphy/internal/scan"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/module"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var versionFlag bool
var consoleLogLevelOverride string

func preRun(cmd *cobra.Command, args []string) error {
	if versionFlag {
		fmt.Println(version.FullInfo())
		fmt.Printf("Supported modules: %s\n", strings.Join(module.GetSupportedModuleList(), ", "))
		os.Exit(0)
	}

	logger.LogFileCleanup()

	if e := common.LogLevel.Of(consoleLogLevelOverride); e != nil {
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

	// version
	c.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "show version and exit")

	// Logging
	c.PersistentFlags().BoolVar(&common.NoLogFile, "no-log-file", false, "do not write log file")
	c.PersistentFlags().StringVar(&common.LogFileOverride, "write-log-to", "", "specify log file path")
	c.PersistentFlags().StringVar(&consoleLogLevelOverride, "log-level", "silent", "specify log level, must be silent|error|warn|info|debug")
	c.PersistentFlags().BoolVar(&common.EnableNetworkLogging, "network-log", false, "print network data")
	_ = c.PersistentFlags().MarkHidden("network-log")

	// API: Authentication & Network
	c.PersistentFlags().StringVar(&common.CliTokenOverride, "token", "", "specify API token")
	c.PersistentFlags().StringVar(&common.CliServerAddressOverride, "server", "", "specify server address")
	c.PersistentFlags().BoolVarP(&env.CliTlsAllowInsecure, "allow-insecure", "x", false, "Allow insecure TLS connection")
	c.PersistentFlags().BoolVar(&env.NoWait, "no-wait", false, "do not wait scan result")

	c.AddCommand(auth.Cmd())
	c.AddCommand(scan.Cmd())
	c.AddCommand(scan.IdeaScan())
	c.AddCommand(scan.DfCmd())
	c.AddCommand(scan.SbomScan())
	c.AddCommand(binscan.Cmd())
	c.AddCommand(internalcmd.Cmd())
	c.AddCommand(internalcmd.MachineIdCmd())
	return c
}

func Execute() error {
	return rootCmd().Execute()
}
