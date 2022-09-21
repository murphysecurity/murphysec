package module

import "github.com/murphysecurity/murphysec/module/rebar3"

func init() {
	Inspectors = append(Inspectors, &rebar3.Inspector{})
}
