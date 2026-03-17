package yarn

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/iseki0/go-yarnlock"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/module/pkgjs"
	"github.com/pkg/errors"
)

func readModuleName(dir string) (string, string) {
	f, e := pkgjs.ReadDir(dir)
	if e != nil {
		return "", ""
	}
	return f.Name, f.Version
}

func yarnFallback(dir string) ([]Dep, error) {
	var rs []Dep
	pkg0, e := pkgjs.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	distinct := map[string]string{}
	for k, v := range pkg0.DevDependencies {
		distinct[k] = v
	}
	for k, v := range pkg0.Dependencies {
		distinct[k] = v
	}

	for k, v := range distinct {
		rs = append(rs, Dep{Name: k, Version: v})
	}
	return rs, nil
}

func analyzeYarnDep(ctx context.Context, dir string) (r []Dep, e error) {
	var logger = logctx.Use(ctx).Sugar()
	f, e := os.Open(filepath.Join(dir, "yarn.lock"))
	if e != nil {
		logger.Infof("Open yarn.lock failed. %v", e)
		return yarnFallback(dir)
	}
	defer func() { _ = f.Close() }()
	data, e := io.ReadAll(io.LimitReader(f, 16*1024*1024))
	if e != nil {
		return nil, errors.Wrap(e, "Read yarn.lock failed.")
	}

	var lockfile yarnlock.LockFile
	var berryData *berryLockData
	newLock := strings.Contains(string(data), "__metadata:")
	if newLock {
		berryData, e = parseYarnLockYamlWithIndex(string(data))
		if e == nil {
			lockfile = berryData.Lockfile
			delete(lockfile, "__metadata")
		}
	} else {
		lockfile, e = yarnlock.ParseLockFileData(data)
	}
	if e != nil {
		return nil, errors.Wrap(e, "Parse lockfile failed.")
	}
	pkg, e := pkgjs.ReadDir(dir)
	if e != nil {
		return nil, e
	}
	if newLock {
		for name, ver := range pkg.Dependencies {
			pkg.Dependencies[name] = "npm:" + ver
		}
		for name, ver := range pkg.DevDependencies {
			pkg.DevDependencies[name] = "npm:" + ver
		}
		return buildDepTreeBerry(berryData, pkg), nil
	}
	return buildDepTree(lockfile, pkg), nil
}
