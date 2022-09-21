package module

import "github.com/murphysecurity/murphysec/module/cargo"

func init() {
	Inspectors = append(Inspectors, &cargo.Inspector{})
}
