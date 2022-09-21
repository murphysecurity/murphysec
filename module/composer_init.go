package module

import "github.com/murphysecurity/murphysec/module/composer"

func init() {
	Inspectors = append(Inspectors, &composer.Inspector{})
}
