package plugin

import (
	"murphysec-cli-simple/plugin/plugin_base"
)

var pluginMap = func() map[string]plugin_base.Plugin {
	m := map[string]plugin_base.Plugin{}
	for i := range Plugins {
		m[Plugins[i].Info().Name] = Plugins[i]
	}
	return m
}()

func GetPluginOrNil(name string) plugin_base.Plugin {
	return pluginMap[name]
}
