//go:build !embedding

package plugin

import (
	"murphysec-cli-simple/plugin/gradle"
	"murphysec-cli-simple/plugin/maven"
	"murphysec-cli-simple/plugin/plugin_base"
)

var Plugins = []plugin_base.Plugin{
	maven.Instance,
	gradle.Instance,
}
