package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/conf"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/module"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var versionFlag bool
var CliServerAddressOverride string
var allowInsecure bool

func rootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:               "murphysec",
		PersistentPreRunE: preRun,
		TraverseChildren:  true,
		Run: func(cmd *cobra.Command, args []string) {
			must.Must(cmd.Help())
		},
	}
	c.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "show version and exit")
	c.PersistentFlags().BoolVar(&disableLogFile, "no-log-file", false, "do not write log file")
	c.PersistentFlags().StringVar(&cliLogFilePathOverride, "write-log-to", "", "specify log file path")
	c.PersistentFlags().StringVar(&consoleLogLevelOverride, "log-level", "silent", "specify log level, must be silent|error|warn|info|debug")
	c.PersistentFlags().BoolVar(&enableNetworkLog, "network-log", false, "print network data")
	c.PersistentFlags().StringVar(&conf.APITokenCliOverride, "token", "", "specify API token")
	c.PersistentFlags().StringVar(&CliServerAddressOverride, "server", "", "specify server address")
	c.PersistentFlags().BoolVarP(&allowInsecure, "allow-insecure", "x", false, "Allow insecure TLS connection")
	c.PersistentFlags().String("ide", "", "hidden")
	must.Must(c.PersistentFlags().MarkHidden("ide"))
	c.AddCommand(authCmd())
	c.AddCommand(scanCmd())
	if env.AllowBinScan {
		c.AddCommand(binScanCmd())
		c.AddCommand(iotScanCmd())
	}
	c.AddCommand(ideaScanCmd())
	if env.ScannerScan {
		c.AddCommand(scannerCmd())
	}
	c.AddCommand(machineCmd())
	c.AddCommand(dockerScanCmd())
	return c
}

func preRun(cmd *cobra.Command, args []string) error {
	if versionFlag {
		fmt.Printf("%s %s\n\n", filepath.Base(must.A(filepath.EvalSymlinks(must.A(os.Executable())))), version.Version())
		var suffix string
		if version.GetGitModified() == "true" {
			suffix = "(modified)"
		}
		fmt.Printf("Git hash: %s%s\n", version.GetGitHash(), suffix)
		fmt.Printf("Git time: %s\n", version.GetGitTime())
		fmt.Printf("Supported modules: %s\n", strings.Join(module.GetSupportedModuleList(), ", "))
		os.Exit(0)
	}
	if !utils.InStringSlice([]string{"", "warn", "error", "debug", "info", "silent"}, strings.ToLower(strings.TrimSpace(consoleLogLevelOverride))) {
		return errors.New("Loglevel invalid, must be silent|error|warn|info|debug")
	}
	if CliServerAddressOverride == "" {
		CliServerAddressOverride = os.Getenv("MPS_CLI_SERVER")
	}
	if CliServerAddressOverride == "" {
		CliServerAddressOverride = "https://www.murphysec.com"
	}
	if allowInsecure {
		// config default http transport allow insecure
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	env.ConfigureServerBaseUrl(CliServerAddressOverride)
	api.C = api.NewClient()
	logFileCleanup()
	api.C.Token = conf.APIToken()
	return nil
}

func Execute() error {
	return rootCmd().Execute()
}
