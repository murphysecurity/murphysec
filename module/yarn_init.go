package module

import "github.com/murphysecurity/murphysec/module/yarn"

func init() {
	Inspectors = append(Inspectors, &yarn.Inspector{})
}
