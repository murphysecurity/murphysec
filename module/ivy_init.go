package module

import "github.com/murphysecurity/murphysec/module/ivy"

func init() {
	Inspectors = append(Inspectors, &ivy.Inspector{})
}
