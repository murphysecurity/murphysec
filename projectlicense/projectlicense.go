package projectlicense

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/murphysecurity/licensematcher"
	"github.com/murphysecurity/murphysec/model"
)

func ScanDir(ctx context.Context) (e error) {
	task := model.UseScanTask(ctx)
	var dir = task.ProjectPath
	entries, e := os.ReadDir(dir)
	if e != nil {
		return
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		var name = entry.Name()
		if strings.Contains(strings.ToLower(name), "license") {
			var path = filepath.Join(dir, name)
			var data []byte
			data, e = os.ReadFile(path)
			if e != nil {
				continue
			}
			var lic = licensematcher.MatchInput(string(data))
			if lic != "" {
				task.ProjectLicense.License = lic
				task.ProjectLicense.Path = name
				return
			}
		}
	}
	return
}
