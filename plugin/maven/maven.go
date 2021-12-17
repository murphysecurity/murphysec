package maven

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/plugin/maven/pom_scanner"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/util/simplejson"
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
	rs, e := doScan(dir)
	if e == nil {
		return nil, e
	}
	output.Warn("Maven execution failed, use another analyzer")
	analyze, err := pom_scanner.Analyze(dir)
	if err != nil {
		return nil, err
	}
	rs = &plugin_base.PackageInfo{
		PackageManager:  "maven",
		PackageFile:     dir,
		PackageFilePath: filepath.Join(dir, "pom.xml"),
		Language:        "java",
		Dependencies:    simplejson.NewFrom(analyze),
		Name:            analyze.Name,
		RuntimeInfo:     simplejson.New(),
	}
	return rs, nil
}

func (p *Plugin) SetupScanCmd(c *cobra.Command) {
}
