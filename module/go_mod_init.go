package module

import "github.com/murphysecurity/murphysec/module/go_mod"

func init() {
	Inspectors = append(Inspectors, &go_mod.Inspector{})
}
