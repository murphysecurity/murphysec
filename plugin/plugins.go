package plugin

import (
	"murphysec-cli-simple/plugin/gradle"
	"murphysec-cli-simple/plugin/hello"
	"murphysec-cli-simple/plugin/plugin_base"
)

var Plugins = []plugin_base.Plugin{
	hello.Instance,
	gradle.Instance,
}
