package composer

import (
	"context"
	"encoding/json"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/simplejson"
	"go.uber.org/zap"
)

func readComposerLockFile(path string) ([]Package, error) {
	lockFileData, e := utils.ReadFileLimited(path, _ComposerLockFileSizeLimit)
	if e != nil {
		return nil, errors.Wrap(e, "Read composer.lock failed")
	}
	pkgs, e := parseComposerLock(lockFileData)
	if e != nil {
		return nil, errors.Wrap(e, "Parse composer.lock failed")
	}
	return pkgs, nil
}

func readComposerInstalledFile(path string) ([]Package, error) {
	lockFileData, e := utils.ReadFileLimited(path, _ComposerLockFileSizeLimit)
	if e != nil {
		return nil, errors.Wrap(e, "Read installed.json failed")
	}
	pkgs, e := parseComposerInstalled(lockFileData)
	if e != nil {
		return nil, errors.Wrap(e, "Parse installed.json failed")
	}
	return pkgs, nil
}

func parseComposerLock(data []byte) ([]Package, error) {
	var j simplejson.JSON
	if e := json.Unmarshal(data, &j); e != nil {
		return nil, errors.Wrap(e, "ParseComposerLock")
	}
	pkgList := make([]Package, 0)
	for _, pkg := range j.Get("packages").JSONArray() {
		p := Package{}
		p.Name = pkg.Get("name").String()
		p.Version = pkg.Get("version").String()
		if p.Name == "" || p.Version == "" {
			continue
		}
		for s := range pkg.Get("require").JSONMap() {
			p.Require = append(p.Require, s)
		}
		pkgList = append(pkgList, p)
	}
	return pkgList, nil
}

func parseComposerInstalled(data []byte) ([]Package, error) {
	var j simplejson.JSON
	if e := json.Unmarshal(data, &j); e != nil {
		return nil, errors.Wrap(e, "ParseComposerInstalled")
	}

	// Composer installed.json has two common schemas:
	// 1) {"packages":[...], ...}
	// 2) [{...}, {...}]
	var packageNodes []*simplejson.JSON
	if arr := j.Get("packages").JSONArray(); len(arr) > 0 {
		packageNodes = arr
	} else if arr := j.JSONArray(); len(arr) > 0 {
		packageNodes = arr
	}

	pkgList := make([]Package, 0, len(packageNodes))
	for _, pkg := range packageNodes {
		if pkg == nil {
			continue
		}
		p := Package{}
		p.Name = pkg.Get("name").String()
		p.Version = pkg.Get("version").String()
		if p.Name == "" || p.Version == "" {
			continue
		}
		for s := range pkg.Get("require").JSONMap() {
			p.Require = append(p.Require, s)
		}
		pkgList = append(pkgList, p)
	}
	return pkgList, nil
}

func readManifest(ctx context.Context, path string) (*Manifest, error) {
	logger := logctx.Use(ctx)
	logger.Debug("readManifest", zap.String("path", path))
	composerFileData, e := utils.ReadFileLimited(path, _ComposerManifestFileSizeLimit)
	if e != nil {
		return nil, wrapErr(ErrReadComposerManifest, e)
	}
	manifest, e := parseComposeManifest(composerFileData)
	if e != nil {
		return nil, wrapErr(ErrParseComposerManifest, e)
	}
	return manifest, nil
}

func parseComposeManifest(data []byte) (*Manifest, error) {
	var j simplejson.JSON
	if e := json.Unmarshal(data, &j); e != nil {
		return nil, e
	}
	m := &Manifest{}
	m.Name = j.Get("name").String()
	m.Version = j.Get("version").String()
	for name, versionConstraint := range j.Get("require").JSONMap() {
		m.Require = append(m.Require, Element{
			Name:    name,
			Version: versionConstraint.String(),
		})
	}
	return m, nil
}
