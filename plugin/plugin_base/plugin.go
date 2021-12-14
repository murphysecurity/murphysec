package plugin_base

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/util/simplejson"
)

type Plugin interface {
	// Info returns PluginInfo of the plugin
	Info() PluginInfo

	// MatchPath returns a boolean indicating if the path is acceptable by the plugin
	MatchPath(dir string) bool

	DoScan(dir string) (*PackageInfo, error)

	SetupScanCmd(c *cobra.Command)
}

type PluginInfo struct {
	// the unique identifier of the plugin
	Name string
	// one-line description.
	ShortDescription string
	// plugin version
	Version string
}

type PackageInfo struct {
	PackageManager  string           `json:"package_manager"`
	PackageFile     string           `json:"package_file"`
	PackageFilePath string           `json:"package_file_path"`
	Language        string           `json:"language"`
	Dependencies    *simplejson.JSON `json:"dependencies"`
	Name            string           `json:"name"`
	RuntimeInfo     *simplejson.JSON `json:"runtime_info"`
}
