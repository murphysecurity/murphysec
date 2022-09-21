package module

import "github.com/murphysecurity/murphysec/module/python"

func init() {
	Inspectors = append(Inspectors, &python.Inspector{})
}
