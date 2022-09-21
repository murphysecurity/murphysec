package module

import "github.com/murphysecurity/murphysec/module/cocoapods"

func init() {
	Inspectors = append(Inspectors, &cocoapods.Inspector{})
}
