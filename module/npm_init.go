package module

import "github.com/murphysecurity/murphysec/module/npm"

func init() {
	Inspectors = append(Inspectors, &npm.Inspector{})
}
