package gradle

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/util/output"
	"os"
	"os/signal"
	"syscall"
)

var Instance plugin_base.Plugin = &Plugin{}

type Plugin struct {
}

func (_ *Plugin) Info() plugin_base.PluginInfo {
	return plugin_base.PluginInfo{Name: "gradle", ShortDescription: "for gradle project"}
}

func (p *Plugin) MatchPath(dir string) bool {
	output.Debug(fmt.Sprintf("gradle - MatchPath: %s", dir))
	f := detectGradleFile(dir)
	if f == "" {
		output.Info("Gradle not detected!")
		return false
	}
	output.Info("Gradle detected!")
	return true
}

func (p *Plugin) DoScan(dir string) (*plugin_base.PackageInfo, error) {
	sigTerm := make(chan os.Signal, 1)
	finishCh := make(chan struct{})
	defer close(finishCh)
	cancel := make(chan struct{})
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		select {
		case <-finishCh:
			signal.Stop(sigTerm)
		case <-sigTerm:
			close(cancel)
		}
	}()

	// do scan
	scanResult, err := scanDir(cancel, dir)
	if err != nil {
		output.Error(fmt.Sprintf("Scan failed, %s", err.Error()))
		return nil, nil
	}
	fmt.Println(scanResult)
	return nil, nil
}

func (p *Plugin) SetupScanCmd(c *cobra.Command) {}
