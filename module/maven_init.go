package module

import "github.com/murphysecurity/murphysec/module/maven"

func init() {
	Inspectors = append(Inspectors, maven.Instance)
}
