package inspector

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
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
	q := list.New()
	q.PushBack(dir)
	for q.Len() > 0 {
		cur := q.Front().Value.(string)
		q.Remove(q.Front())
		if entries, e := os.ReadDir(cur); e != nil {
			continue
		} else {
			for _, entry := range entries {
				name := entry.Name()
				// check name black list
				if checkNameBlackList(name) {
					continue
				}

				p := filepath.Join(cur, name)
				info, e := entry.Info()
				if e != nil {
					logger.Info.Println("Get info failed", p)
					continue
				}

				// check file mode
				fileMod := info.Mode()
				switch fileMod {
				case os.ModeCharDevice, os.ModeDevice, os.ModeNamedPipe, os.ModeSocket, os.ModeSymlink, os.ModeIrregular:
					continue
				}

				if entry.IsDir() {
					q.PushBack(p)
					continue
				}
				// check file size
				if info.Size() < 32 {
					continue
				}

				pathCh <- p
			}
		}
	}
}

func checkNameBlackList(name string) bool {
	return name == "node_modules" || strings.HasPrefix(name, ".")
}
