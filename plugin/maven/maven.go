package maven

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/must"
	"path/filepath"
)

type Plugin struct {
}

var Instance plugin_base.Plugin = &Plugin{}

func (_ *Plugin) Info() *plugin_base.PluginInfo {
	return &plugin_base.PluginInfo{Name: "maven", ShortDescription: "for maven package"}
}

func (p *Plugin) MatchPath(dir string) bool {
	f := filepath.Join(must.String(filepath.Abs(dir)), "pom.xml")
	return util.IsPathExist(f) && !util.IsDir(f)
}

func (p *Plugin) DoScan(dir string) (*plugin_base.PackageInfo, error) {
	return doScan(dir)
}

func (p *Plugin) SetupScanCmd(c *cobra.Command) {
}
