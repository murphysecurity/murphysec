package plugin_base

import (
	"github.com/spf13/cobra"
)

type Plugin interface {
	// Info returns PluginInfo of the plugin
	Info() PluginInfo

	// MatchPath returns a boolean indicating if the path is acceptable by the plugin
	MatchPath(dir string) bool

	DoScan(dir string) interface{}

	SetupScanCmd(c *cobra.Command)
}

type PluginInfo struct {
	// the unique identifier of the plugin
	Name string
	// one-line description.
	ShortDescription string
}
