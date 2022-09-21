package module

import "github.com/murphysecurity/murphysec/module/conan"

func init() {
	Inspectors = append(Inspectors, &conan.Inspector{})
}
