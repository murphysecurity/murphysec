package gradle

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/scanner"
	"murphysec-cli-simple/util/output"
	"os"
	"path/filepath"
)

var Instance Plugin

type Plugin struct {
}

func (_ *Plugin) Info() scanner.PluginInfo {
	return scanner.PluginInfo{Name: "gradle", ShortDescription: "for gradle project"}
}

func (p *Plugin) MatchPath(ctx context.Context, dir string) bool {
	output.Debug(fmt.Sprintf("gradle - MatchPath: %s", dir))
	f := detectGradleFile(dir)
	if f == "" {
		output.Info("Gradle not detected!")
		return false
	}
	output.Info("Gradle detected!")
	return true
}

func (p *Plugin) DoScan(ctx context.Context, dir string) interface{} {
	// todo: scan
	panic("todo")
	return nil
}

func (p *Plugin) SetupScanCmd(c *cobra.Command) {}

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
