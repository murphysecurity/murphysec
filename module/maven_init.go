package module

import (
	"github.com/murphysecurity/murphysec/module/maven"
	"os"
)

func init() {
	if os.Getenv("DO_NOT_SCAN_MAVEN") == "1" {
		return
	}
	Inspectors = append(Inspectors, &maven.Inspector{})
}
