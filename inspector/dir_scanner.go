package inspector

import (
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"os"
	"path/filepath"
)

type dirScanner struct {
	inspectors  []model.Inspector
	scannedDirs []dirScanItem
	root        string
}

type dirScanItem struct {
	path      string
	inspector model.Inspector
}

func (d *dirScanItem) String() string {
	return fmt.Sprintf("%s - %s", d.inspector, d.path)
}

func (d *dirScanner) scan() {
	d._r(0, d.root, map[model.Inspector]unit{})
}

func (d *dirScanner) _r(depth int, p string, usedInspector map[model.Inspector]unit) {
	if depth > 16 {
		return
	}
	entries, e := os.ReadDir(p)
	if e != nil {
		return
	}

	for _, it := range d.inspectors {
		_, used := usedInspector[it]
		if used {
			if !it.SupportFeature(model.InspectorFeatureAllowNested) {
				continue
			}
			if it.CheckDir(p) {
				d.scannedDirs = append(d.scannedDirs, dirScanItem{
					path:      p,
					inspector: it,
				})
			}
		} else {
			if it.CheckDir(p) {
				usedInspector[it] = unit{}
				// Clear the used flag of the first time using
				//goland:noinspection ALL
				defer delete(usedInspector, it)
				d.scannedDirs = append(d.scannedDirs, dirScanItem{
					path:      p,
					inspector: it,
				})
			}
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		entryName := entry.Name()
		if dirShouldIgnore(entryName) {
			continue
		}
		d._r(depth+1, filepath.Join(p, entryName), usedInspector)
	}
}
