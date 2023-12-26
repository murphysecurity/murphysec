package module

import "github.com/murphysecurity/murphysec/module/arkts"

func init() {
	Inspectors = append(Inspectors, arkts.Inspector{})
}
