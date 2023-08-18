package module

import (
	"github.com/murphysecurity/murphysec/module/gradle"
	"os"
)

func init() {
	if os.Getenv("DO_NOT_SCAN_GRADLE") == "1" {
		return
	}
	Inspectors = append(Inspectors, &gradle.Inspector{})
}
