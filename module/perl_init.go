package module

import "github.com/murphysecurity/murphysec/module/perl"

func init() {
	Inspectors = append(Inspectors, &perl.Inspector{})
}
