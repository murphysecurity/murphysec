package inspector

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

const _FileHashScanConcurrency = 2

func FileHashScan(ctx context.Context) error {
	logger := Logger.Named("file-hash-scan")
	task := model.UseScanTask(ctx)
	basePath, e := filepath.Abs(task.ProjectDir)
	if e != nil {
		return errors.Wrap(e, "Get absolute path fail.")
	}
	logger.Info("File hash scan", zap.String("basePath", basePath))

	filepathCh, findCxxFileRoutine := goFindAllCxxFile(basePath, logger.Named("cxx-file-finder"))
	go findCxxFileRoutine(ctx)
	fileHashCh, calcHashRoutine := goCalcFileHash(filepathCh, _FileHashScanConcurrency, logger.Named("calc-hash"))
	go calcHashRoutine(ctx)

	for it := range fileHashCh {
		p, e := filepath.Rel(basePath, it.Path)
		if e != nil {
			logger.Error("Get relative-path fail, skip.", zap.Error(e))
			continue
		}
		it.Path = p
		task.FileHashes = append(task.FileHashes, it)
	}
	return nil
}

func goCalcFileHash(filepathInputCh chan string, concurrency int, logger *zap.Logger) (fileHashCh chan model.FileHash, routine func(ctx context.Context)) {
	if logger == nil {
		logger = zap.NewNop()
	}
	fileHashCh = make(chan model.FileHash, 16)
	routine = func(ctx context.Context) {
		defer close(fileHashCh)
		wg := sync.WaitGroup{}
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case fp, ok := <-filepathInputCh:
						if !ok {
							return
						}
						s, e := calcFileHash(fp)
						if e != nil {
							logger.Warn("Calc file hash error", zap.Error(e), zap.String("path", fp))
							continue
						}
						select {
						case fileHashCh <- model.FileHash{Path: fp, Hash: s}:
						case <-ctx.Done():
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		wg.Wait()
	}
	return
}

func goFindAllCxxFile(baseDir string, logger *zap.Logger) (filepathCh chan string, routine func(ctx context.Context)) {
	filepathCh = make(chan string, 16)
	if logger == nil {
		logger = zap.NewNop()
	}
	if !filepath.IsAbs(baseDir) {
		panic("base dir must be absolute")
	}
	routine = func(ctx context.Context) {
		defer close(filepathCh)
		if ctx == nil {
			ctx = context.TODO()
		}
		var counter int
		logger.Debug("Start walker", zap.String("dir", baseDir))
		walkDirFunc := func(path string, d fs.DirEntry, err error) error {
			if ctx.Err() != nil {
				return nil
			}
			if d == nil || err != nil {
				logger.Info("Walker error", zap.String("path", path), zap.Error(err))
				return nil
			}
			if dirShouldIgnore(path) {
				if d.IsDir() {
					return fs.SkipDir
				} else {
					return nil
				}
			}
			if _CxxExtSet[filepath.Ext(d.Name())] {
				select {
				case <-ctx.Done():
				case filepathCh <- path:
					counter++
				}
			}
			return nil
		}
		e := filepath.WalkDir(baseDir, walkDirFunc)
		logger.Debug("Walker terminated", zap.String("dir", baseDir), zap.Int("total", counter))
		if e != nil {
			logger.Warn("Walker error", zap.Error(e))
		}
	}
	return
}

func calcFileHash(path string) ([]string, error) {
	var rs []string
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	h1 := md5.New()
	h2 := md5.New()
	h3 := md5.New()
	w1 := h1
	w2 := utils.Dos2UnixWriter(h2)
	w3 := utils.Unix2DosWriter(h3)
	w := io.MultiWriter(w1, w2, w3)
	if _, e := io.Copy(w, f); e != nil {
		return nil, e
	}
	_ = w2.Close()
	_ = w3.Close()
	rs = append(rs, hex.EncodeToString(h1.Sum(make([]byte, 0, 16))))
	rs = append(rs, hex.EncodeToString(h2.Sum(make([]byte, 0, 16))))
	rs = append(rs, hex.EncodeToString(h3.Sum(make([]byte, 0, 16))))

	return utils.DistinctStringSlice(rs), nil
}

var _CxxExtSet = map[string]bool{
	".c":   true,
	".h":   true,
	".hpp": true,
	".cpp": true,
	".cxx": true,
	".c++": true,
	".C":   true,
	".cc":  true,
	".hxx": true,
	".C++": true,
	".cp":  true,
}
