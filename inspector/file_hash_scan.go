package inspector

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

	filepathCh := make(chan string, 32)
	go func() {
		findAllCxxFile(basePath, filepathCh)
		defer close(filepathCh)
	}()

	// file hash
	hashChan := make(chan model.FileHash, 32)
	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < _FileHashScanConcurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mapFilepathToHash(filepathCh, hashChan)
			}()
		}
		wg.Wait()
		close(hashChan)
	}()

	for it := range hashChan {
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

func mapFilepathToHash(fileCh chan string, outputCh chan model.FileHash) {
	for file := range fileCh {
		hash := calcFileHashIgnoreErr(file)
		outputCh <- model.FileHash{
			Path: file,
			Hash: hash,
		}
	}
}

func findAllCxxFile(baseDir string, filepathCh chan string) {
	Logger.Info("findAllCxxFile", zap.String("baseDir", baseDir))
	filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			Logger.Error("DirEntry is nil, skip", zap.Error(err))
			return nil
		}
		if checkNameBlackList(d.Name()) {
			if d.IsDir() {
				return fs.SkipDir
			} else {
				return nil
			}
		}
		if CxxExtSet[filepath.Ext(d.Name())] {
			filepathCh <- path
		}
		return nil
	})
}

func calcFileHashIgnoreErr(path string) []string {
	s, e := calcFileHash(path)
	if e != nil {
		return nil
	}
	return s
}

func calcFileHash(path string) ([]string, error) {
	var rs []string
	f, e := os.Open(path)
	if e != nil {
		return nil, errors.Wrap(e, fmt.Sprintf("Open file failed when calc file hash: %s", path))
	}
	defer f.Close()

	// TODO: switch to SHA-256
	h1 := md5.New()
	h2 := md5.New()
	h3 := md5.New()
	w1 := h1
	w2 := utils.Dos2UnixWriter(h2)
	w3 := utils.Unix2DosWriter(h3)
	w := io.MultiWriter(w1, w2, w3)
	if _, e := io.Copy(w, f); e != nil {
		return nil, errors.Wrap(e, fmt.Sprintf("Calc file hash failed %s", path))
	}
	_ = w2.Close()
	_ = w3.Close()
	rs = append(rs, hex.EncodeToString(h1.Sum(make([]byte, 0, 16))))
	rs = append(rs, hex.EncodeToString(h2.Sum(make([]byte, 0, 16))))
	rs = append(rs, hex.EncodeToString(h3.Sum(make([]byte, 0, 16))))

	return utils.DistinctStringSlice(rs), nil
}

func checkNameBlackList(name string) bool {
	return name == "node_modules" || strings.HasPrefix(name, ".")
}

var CxxExtSet = map[string]bool{
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
