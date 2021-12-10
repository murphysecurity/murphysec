package plugins

import (
	"murphysec-cli-simple/plugins/hello"
	"murphysec-cli-simple/scanner"
)

var Plugins = []scanner.Plugin{
	&hello.Instance,
}
