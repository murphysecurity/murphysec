package module

import "github.com/murphysecurity/murphysec/module/gradle"

func init() {
	Inspectors = append(Inspectors, &gradle.Inspector{})
}
