package inspector

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

func managedInspect(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	baseDir := scanTask.ProjectDir
	Logger.Info("Auto scan dir", zap.String("dir", baseDir))

	var scanner = &dirScanner{
		inspectors: module.Inspectors,
		root:       baseDir,
	}
	scanner.scan()

	Logger.Sugar().Infof("Found %d directories, in %v", len(scanner.scannedDirs), time.Now().Sub(scanTask.StartTime))
	for idx, it := range scanner.scannedDirs {
		st := time.Now()
		c := model.WithInspectorTask(ctx, it.path)
		c = utils.WithLogger(c, Logger.Named(fmt.Sprintf("%s-%d", it.inspector.String(), idx)))
		Logger.Sugar().Infof("Begin: %s, duration: %v", it.String(), time.Now().Sub(st))
		e := it.inspector.InspectProject(c)
		Logger.Sugar().Infof("End: %s, duration: %v", it.String(), time.Now().Sub(st))
		if e != nil {
			Logger.Error("InspectError", zap.Error(e), zap.Any("inspector", it))
		}
	}
	return nil
}

type dirScanner struct {
	inspectors  []base.Inspector
	scannedDirs []dirScanItem
	root        string
}

type dirScanItem struct {
	path      string
	inspector base.Inspector
}

func (d *dirScanItem) String() string {
	return fmt.Sprintf("%s - %s", d.inspector, d.path)
}

func (d *dirScanner) scan() {
	d._r(0, d.root, map[base.Inspector]unit{})
}

func (d *dirScanner) _r(depth int, p string, usedInspector map[base.Inspector]unit) {
	if depth > 6 {
		return
	}
	entries, e := os.ReadDir(p)
	if e != nil {
		return
	}

	for _, it := range d.inspectors {
		_, used := usedInspector[it]
		if used {
			if !it.SupportFeature(base.InspectorFeatureAllowNested) {
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
