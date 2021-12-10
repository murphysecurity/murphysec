package plugin_base

import (
	"context"
	"github.com/spf13/cobra"
)

type Plugin interface {
	// Info returns PluginInfo of the plugin
	Info() PluginInfo

	// MatchPath returns a boolean indicating if the path is acceptable by the plugin
	MatchPath(ctx context.Context, dir string) bool

	DoScan(ctx context.Context, dir string) interface{}

	SetupScanCmd(c *cobra.Command)
}

type PluginInfo struct {
	// the unique identifier of the plugin
	Name string
	// one-line description.
	ShortDescription string
}
