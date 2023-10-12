package pnpm

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/pnpm/shared"
	v5 "github.com/murphysecurity/murphysec/module/pnpm/v5"
	"io"
	"os"
	"path/filepath"
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
	version, e := parseLockfileVersion(data)
	if e != nil {
		result.e = fmt.Errorf("parse lockfile version failed, %w", e)
		return
	}
	versionNumber := matchLockfileVersion(version)
	if versionNumber == 5 {
		lockfile, e := v5.ParseLockfile(data)
		if e != nil {
			result.e = fmt.Errorf("v5: %w", e)
			return
		}
		result.trees = v5.AnalyzeDepTree(lockfile)
	} else if versionNumber == 6 {
		// todo: v6 support need rewrite
		lockfile, e := parseV6Lockfile(data, false)
		if e != nil {
			result.e = fmt.Errorf("v6: %w", e)
			return
		}
		items, e := lockfile.buildDependencyTree(false)
		if e != nil {
			result.e = fmt.Errorf("v6: %w", e)
			return
		}
		result.trees = []shared.DepTree{{
			Name:         "",
			Dependencies: items,
		}}
	} else {
		result.e = fmt.Errorf("unsupported version \"%s\"", version)
		return
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
