package module

import "github.com/murphysecurity/murphysec/module/poetry"

func init() {
	Inspectors = append(Inspectors, &poetry.Inspector{})
}
