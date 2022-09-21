package module

import "github.com/murphysecurity/murphysec/module/sbt"

func init() {
	Inspectors = append(Inspectors, &sbt.Inspector{})
}
