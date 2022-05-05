package inspector

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const _FileHashScanConcurrency = 2

func FileHashScan(ctx *ScanContext) error {
	basePath, e := filepath.Abs(ctx.ProjectDir)
	if e != nil {
		return errors.Wrap(e, "Get absolute path fail.")
	}
	logger.Info.Println("FileHashScan: ", basePath)

	filepathCh := make(chan string, 32)
	go func() {
		findAllCxxFile(basePath, filepathCh)
		defer close(filepathCh)
	}()

	// file hash
	hashChan := make(chan FileHash, 32)
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
			logger.Err.Println("Get relative-path fail, skip.", e.Error())
			continue
		}
		it.Path = p
		ctx.FileHashes = append(ctx.FileHashes, it)
	}
	return nil
}

func mapFilepathToHash(fileCh chan string, outputCh chan FileHash) {
	for file := range fileCh {
		hash := calcFileHashIgnoreErr(file)
		outputCh <- FileHash{
			Path: file,
			Hash: hash,
		}
	}
}

func findAllCxxFile(baseDir string, filepathCh chan string) {
	logger.Info.Println("findAllCxxFile: ", baseDir)
	filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Warn.Println("WalkDir error:", err.Error())
			return nil
		}
		if d == nil {
			logger.Warn.Println("DirEntry is nil, skip")
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
