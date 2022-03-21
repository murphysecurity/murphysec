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

func FileHashScan(ctx *ScanContext) {
	fileCh := make(chan string, 16)
	m := sync.Mutex{}
	fileHashes := map[string]interface{}{}
	go dirScan(ctx.ProjectDir, fileCh)
	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			for p := range fileCh {
				for _, it := range calcFileHashIgnoreErr(p) {
					m.Lock()
					fileHashes[it] = struct{}{}
					m.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for s := range fileHashes {
		ctx.FileHashes = append(ctx.FileHashes, s)
	}
	logger.Info.Println("total:", len(fileHashes), "hashes")
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

func dirScan(dir string, pathCh chan string) {
	logger.Info.Printf("dir scan: %s", dir)
	defer logger.Info.Println("dir scan terminated")
	defer close(pathCh)
	e := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if info.Size() < 32 || info.Size() > 16*1024*1024 {
			return nil
		}
		pathCh <- filepath.Join(path)
		return nil
	})
	if e != nil {
		logger.Warn.Println("filepath.Walk err:", e.Error())
	}
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
