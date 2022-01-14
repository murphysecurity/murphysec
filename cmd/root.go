package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"murphysec-cli-simple/version"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var versionFlag bool
var managedMode bool
var logRedirectPath string
var noLogFile bool
var logLevel string

func rootCmd() *cobra.Command {
	argsMap := map[string]bool{}
	for _, it := range os.Args {
		argsMap[it] = true
	}
	c := &cobra.Command{
		Use:               "murphysec",
		PersistentPreRunE: preRun,
		TraverseChildren:  true,
		Run: func(cmd *cobra.Command, args []string) {
			must.Must(cmd.Help())
		},
	}
	c.PersistentFlags().BoolVar(&versionFlag, "version", false, "show version and exit")
	c.PersistentFlags().BoolVar(&noLogFile, "no-log-file", false, "do not write log file")
	c.PersistentFlags().StringVar(&logRedirectPath, "write-log-to", "", "specify log file path")
	c.PersistentFlags().StringVar(&logLevel, "log-level", "", "specify log level")
	c.PersistentFlags().StringVar(&conf.APITokenCliOverride, "token", "", "specify API token")
	c.PersistentFlags().StringVar(&api.CliServerAddressOverride, "server", "", "specify server address")
	// workaround avoid err
	if argsMap["--managed-mode"] {
		c.PersistentFlags().BoolVar(&managedMode, "managed-mode", false, "")
	}
	if argsMap["--ide"] {
		c.PersistentFlags().String("ide", "", "ignore")
	}
	c.AddCommand(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello world!")
			logger.Debug.Println("log...")
		},
	})
	c.AddCommand(authCmd())
	c.AddCommand(scanCmd())
	if argsMap["ideascan"] {
		c.AddCommand(ideaScanCmd())
	}
	return c
}

func preRun(cmd *cobra.Command, args []string) error {
	if versionFlag {
		version.PrintVersionInfo()
		return nil
	}
	if !noLogFile {
		if logRedirectPath == "" {
			logRedirectPath = must.String(homedir.Expand(filepath.Join("~", ".murphysec", "logs", fmt.Sprintf("%d.log", time.Now().UnixMilli()))))
		}
		logger.InitLogFile(logRedirectPath)
	}
	if logLevel != "" {
		switch strings.ToLower(strings.TrimSpace(logLevel)) {
		case "debug":
			logger.SetConsoleLogLevel(logger.LogDebug)
		case "info":
			logger.SetConsoleLogLevel(logger.LogInfo)
		case "warn":
			logger.SetConsoleLogLevel(logger.LogWarn)
		case "error":
			logger.SetConsoleLogLevel(logger.LogErr)
		default:
			return errors.New("Bad log level, must be debug|info|warn|error")
		}
	}
	return nil
}

func Execute() error {
	return rootCmd().Execute()
}
