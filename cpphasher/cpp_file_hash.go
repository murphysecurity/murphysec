package cpphasher

import (
	"context"
	"crypto/md5"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/nocrlfpipe"
	"github.com/murphysecurity/murphysec/utils/must"
	"golang.org/x/sync/errgroup"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func MD5HashingCppFiles(parentCtx context.Context, dir string) (rs [][16]byte, e error) {
	var (
		fileCh  = make(chan string, 1)
		hashCh  = make(chan fileHash16, 1)
		eg, ctx = errgroup.WithContext(parentCtx)
		set     = make(map[[16]byte]struct{})
	)
	// finding files
	eg.Go(func() error { defer close(fileCh); return findAllCppHashingFiles(ctx, dir, fileCh) })
	// hashing files
	eg.Go(func() error {
		defer close(hashCh)
		return md5HashingFilesFromChannelConcurrently(ctx, 2, fileCh, hashCh)
	})

	for hash := range hashCh {
		set[hash.hash] = struct{}{}
	}

	for bytes := range set {
		rs = append(rs, bytes)
	}
	e = eg.Wait()
	return
}

// findAllCppFiles writes abs path which should be hashed to files channel in dir
func findAllCppHashingFiles(ctx context.Context, dir string, files chan<- string) (e error) {
	var (
		logger = logctx.Use(ctx).Sugar()
	)
	must.True(filepath.IsAbs(dir))
	logger.Debugf("finding cpp files: %s", dir)
	var counter int
	defer func() { logger.Debugf("finding cpp files completed, total: %d, %v", counter, e) }() // catch variable counter, ugly golang
	e = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d == nil || err != nil {
			logger.Debugf("walk: %v", err)
			return err
		}
		var fn = d.Name()
		if d.IsDir() {
			if dirShouldIgnore(fn) {
				return filepath.SkipDir
			}
			return nil
		}
		if len(fn) > 0 && fn[0] == '.' {
			return nil
		}
		if !cppFileExtSet[filepath.Ext(fn)] {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case files <- path:
			counter++
			return nil
		}
	})
	return
}

type fileHash16 struct {
	file string
	hash [16]byte
}

func md5HashingFilesFromChannelConcurrently(parentCtx context.Context, concurrency int, files <-chan string, hashes chan<- fileHash16) (e error) {
	var (
		eg, ctx = errgroup.WithContext(parentCtx)
	)
	must.True(concurrency > 0)
	for i := 0; i < concurrency; i++ {
		eg.Go(func() error { return md5HashingFilesFromChannel(ctx, files, hashes) })
	}
	return eg.Wait()
}

func md5HashingFilesFromChannel(ctx context.Context, files <-chan string, hashes chan<- fileHash16) (e error) {
	var (
		logger = logctx.Use(ctx).Sugar()
		hash   [16]byte
	)
	logger.Debugf("md5HashingFilesFromChannel start")
	defer logger.Debugf("md5HashingFilesFromChannel complete, %v", e)
	for {
		select {
		case f, ok := <-files:
			if !ok {
				return
			}
			hash, e = calcFileMD5(f)
			if e != nil {
				logger.Debugf("hash file failed, file: %s, error: %e", f, e)
				continue
			}
			select {
			case hashes <- fileHash16{f, hash}:
			case <-ctx.Done():
				e = ctx.Err()
				return
			}
		case <-ctx.Done():
			e = ctx.Err()
			return
		}
	}
}

func calcFileMD5(path string) (hash [16]byte, e error) {
	var (
		f *os.File
	)

	f, e = os.Open(path)
	if e != nil {
		return
	}
	defer func() { _ = f.Close() }()

	h1 := md5.New()
	w1 := nocrlfpipe.NewNoCrlfWriter(h1)
	must.M(w1.Close())
	_, e = io.Copy(w1, f)
	if e != nil {
		return
	}
	h1.Sum(hash[:0])
	return
}
