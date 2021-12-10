package gradle

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/scanner"
	"murphysec-cli-simple/util/output"
)

var Instance Plugin

type Plugin struct {
}

func (_ *Plugin) Info() scanner.PluginInfo {
	return scanner.PluginInfo{Name: "gradle", ShortDescription: "for gradle project"}
}

func (p *Plugin) MatchPath(ctx context.Context, dir string) bool {
	output.Debug(fmt.Sprintf("gradle - MatchPath: %s", dir))
	f := detectGradleFile(dir)
	if f == "" {
		output.Info("Gradle not detected!")
		return false
	}
	output.Info("Gradle detected!")
	return true
}

func (p *Plugin) DoScan(ctx context.Context, dir string) interface{} {
	// todo: scan
	panic("todo")
	return nil
}

func (p *Plugin) SetupScanCmd(c *cobra.Command) {}
