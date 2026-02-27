package pnpm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
	v5 "github.com/murphysecurity/murphysec/module/pnpm/v5"
	v6 "github.com/murphysecurity/murphysec/module/pnpm/v6"
	v9 "github.com/murphysecurity/murphysec/module/pnpm/v9"
)

var EcoRepo = model.EcoRepo{
	Ecosystem:  "npm",
	Repository: "",
}

const LockfileName = "pnpm-lock.yaml"
const MaxLockfileSize = 32 * 1024 * 1024

type processDirResult struct {
	trees    []shared.DepTree
	lockfile string
	e        error
}

func processDir(ctx context.Context, dir string) (result processDirResult) {
	result.lockfile = filepath.Join(dir, LockfileName)
	LOG := logctx.Use(ctx).Sugar()
	f, e := openLockfile(ctx, result.lockfile)
	if e != nil {
		result.e = e
		return
	}
	LOG.Debugf("reading %s(%s)", LockfileName, result.lockfile)
	data, e := io.ReadAll(io.LimitReader(f, MaxLockfileSize))
	if e != nil {
		result.e = fmt.Errorf("reading %s failed: %w", LockfileName, e)
		return
	}
	processLockfile(ctx, data, &result)
	return
}

func processLockfile(ctx context.Context, data []byte, result *processDirResult) {
	version, e := parseLockfileVersion(data)
	if e != nil {
		result.e = fmt.Errorf("parse lockfile version failed, %w", e)
		return
	}
	versionNumber := matchLockfileVersion(version)
	if versionNumber == 5 {
		result.trees, e = processV5(ctx, data)
		if e != nil {
			result.e = fmt.Errorf("v5: %w", e)
			return
		}
	} else if versionNumber == 6 {
		result.trees, e = v6.Process(ctx, data, false)
		if e != nil {
			result.e = fmt.Errorf("v6: %w", e)
			return
		}
	} else if versionNumber == 9 {
		result.trees, e = v9.Parse(ctx, bytes.NewReader(data))
		if e != nil {
			result.e = fmt.Errorf("v9: %w", e)
			return
		}
	} else {
		result.e = fmt.Errorf("unsupported version \"%s\"", version)
		return
	}
	normalizeTrees(result.trees)
}

func normalizeTrees(trees []shared.DepTree) {
	for i := range trees {
		sortDepItems(trees[i].Dependencies)
	}
	sort.SliceStable(trees, func(i, j int) bool {
		return trees[i].Name < trees[j].Name
	})
}

func sortDepItems(deps []model.DependencyItem) {
	for i := range deps {
		sortDepItems(deps[i].Dependencies)
	}
	sort.SliceStable(deps, func(i, j int) bool {
		if deps[i].CompName != deps[j].CompName {
			return deps[i].CompName < deps[j].CompName
		}
		if deps[i].CompVersion != deps[j].CompVersion {
			return deps[i].CompVersion < deps[j].CompVersion
		}
		if deps[i].IsOnline.Value != deps[j].IsOnline.Value {
			return !deps[i].IsOnline.Value && deps[j].IsOnline.Value
		}
		return deps[i].IsOnline.Valid && !deps[j].IsOnline.Valid
	})
}

func processV5(ctx context.Context, data []byte) (trees []shared.DepTree, e error) {
	lockfile, e := v5.ParseLockfile(data)
	if e != nil {
		return nil, fmt.Errorf("v5: %w", e)
	}
	if r := v5.BuildDepTree(lockfile, nil, ""); r != nil {
		trees = append(trees, *r)
	}
	for s, i := range lockfile.Importers {
		if r := v5.BuildDepTree(lockfile, i, s); r != nil {
			trees = append(trees, *r)
		}
	}
	return
}

func openLockfile(ctx context.Context, lockfilePath string) (f *os.File, e error) {
	f, e = os.Open(lockfilePath)
	if e != nil {
		if os.IsNotExist(e) {
			return nil, lockfileNotExistError{filename: LockfileName, e: e}
		}
		return nil, fmt.Errorf("open lockfile failed(%s): %w", lockfilePath, e)
	}
	defer func() {
		if e != nil {
			closeError := f.Close()
			if closeError != nil {
				logctx.Use(ctx).Sugar().Warnf("closing lockfile %s failed: %s", lockfilePath, e.Error())
			}
		}
	}()
	stat, e := f.Stat()
	if e != nil {
		return nil, fmt.Errorf("retrieve file stat failed(%s): %w", lockfilePath, e)
	}
	if !stat.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file: %w", LockfileName, e)
	}
	if stat.Size() > MaxLockfileSize {
		return nil, fmt.Errorf("%s too big", LockfileName)
	}
	return f, nil
}

type lockfileNotExistError struct {
	filename string
	e        error
}

func (l lockfileNotExistError) Error() string {
	return l.filename + " is not exists: " + l.e.Error()
}

func (l lockfileNotExistError) Unwrap() error {
	return l.e
}
