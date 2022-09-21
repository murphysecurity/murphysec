package module

import "github.com/murphysecurity/murphysec/module/nuget"

func init() {
	Inspectors = append(Inspectors, &nuget.Inspector{})
}
