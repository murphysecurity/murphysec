package module

import "github.com/murphysecurity/murphysec/module/bundler"

func init() {
	Inspectors = append(Inspectors, &bundler.Inspector{})
}
