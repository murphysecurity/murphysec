package gradle

import (
	"fmt"
	"murphysec-cli-simple/util/output"
	"os"
	"path/filepath"
)

// detectGradleFile returns gradle file path in dir, returns nil if not found.
func detectGradleFile(dir string) string {
	for s := range gradleFiles {
		p := filepath.Join(dir, s)
		output.Debug(fmt.Sprintf("try to detect gradle file: %s", p))
		if stat, err := os.Stat(filepath.Join(dir, s)); err == nil && !stat.IsDir() {
			output.Debug("found")
			return p
		}
	}
	output.Debug(fmt.Sprintf("not found any gradle file under: %s", dir))
	return ""
}

var gradleFiles = map[string]bool{
	"build.gradle":     true,
	"build.gradle.kts": true,
}
