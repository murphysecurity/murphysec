package hello

import (
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/plugin/plugin_base"
)

type Plugin struct {
	arg string
}

var Instance plugin_base.Plugin = &Plugin{}

func (_ *Plugin) Info() plugin_base.PluginInfo {
	return plugin_base.PluginInfo{Name: "hello", ShortDescription: "just a hello world"}
}

func (p *Plugin) MatchPath(dir string) bool {
	fmt.Println("hello world MatchPath", p.arg)
	return false
}

func (p *Plugin) DoScan(dir string) (*plugin_base.PackageInfo, error) {
	fmt.Println("hello world DoScan", p.arg)
	return nil, nil
}

func (p *Plugin) SetupScanCmd(c *cobra.Command) {
	c.PersistentFlags().StringVarP(&p.arg, "foo", "", "", "--foo bar")
}
