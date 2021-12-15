package plugin

import (
	"murphysec-cli-simple/plugin/gradle"
	"murphysec-cli-simple/plugin/hello"
	"murphysec-cli-simple/plugin/maven"
	"murphysec-cli-simple/plugin/plugin_base"
)

var Plugins = []plugin_base.Plugin{
	hello.Instance,
	gradle.Instance,
	maven.Instance,
}

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
