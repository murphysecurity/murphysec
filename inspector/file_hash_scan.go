package inspector

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func FileHashInspect(ctx *ScanContext) (*api.VoDetectResponse, error) {
	m, errCount := FileHashInspectScan(ctx.ProjectDir)
	// distinct hash
	hashStr := make(map[string]struct{})
	for _, s := range m {
		for _, it := range s {
			hashStr[hex.EncodeToString(it)] = struct{}{}
		}
	}
	if ctx.TaskSource == api.TaskSourceCli {
		fmt.Println(fmt.Sprintf("共成功扫描 %d 个文件，提交 %d 个哈希，%d 个文件扫描失败", len(m), len(hashStr), errCount))
	}
	// api request object
	req := ctx.getApiRequestObj()
	moduleVo := api.VoModule{
		FileHashList:   []api.FileHash{},
		Language:       "C,C++",
		Name:           ctx.ProjectName,
		PackageFile:    ctx.ProjectName,
		PackageManager: "unmanaged",
		RelativePath:   "/",
		RuntimeInfo:    nil,
		Version:        "",
		ModuleUUID:     uuid.UUID{},
	}
	for s := range hashStr {
		moduleVo.FileHashList = append(moduleVo.FileHashList, api.FileHash{Hash: s})
	}
	req.Modules = []api.VoModule{moduleVo}
	return api.SendDetectHash(req)
}

const ScanConcurrency = 2

// todo: refactor

func FileHashInspectScan(projectDir string) (map[string][][]byte, int) {

	fch := make(chan interface{}, 100)
	go func() {
		_scanRecursive(projectDir, fch, map[string]struct{}{})
		close(fch)
	}()
	wg := sync.WaitGroup{}
	type hashResult struct {
		path string
		hash [][]byte
		err  error
	}
	// 并发扫描
	hashCh := make(chan hashResult, 10)
	for i := 0; i < ScanConcurrency; i++ {
		wg.Add(1)
		go func() {
			for it := range fch {
				switch v := it.(type) {
				case error:
					logger.Warn.Println(v.Error())
				case string:
					if stat, e := os.Stat(v); e != nil {
						logger.Debug.Println("Get file info failed, skip", v, e.Error())
						continue
					} else {
						if stat.Size() < 32 {
							logger.Debug.Println("File size < 32 bytes, skip", v, stat.Size(), "bytes")
							continue
						}
					}
					hash, e := calcFileHash(v)
					hashCh <- hashResult{
						path: v,
						hash: hash,
						err:  e,
					}
				}
			}
			wg.Done()
		}()
	}
	// 收集结果到map
	type result struct {
		errCount int
		result   map[string][][]byte
	}
	resultCh := make(chan result, 1)
	go func() {
		rs := result{
			errCount: 0,
			result:   map[string][][]byte{},
		}
		for it := range hashCh {
			if it.err != nil {
				rs.errCount++
				logger.Warn.Printf("Calc hash failed: %+v\n", it.err)
				continue
			}
			var hashS []string
			for _, it := range it.hash {
				hashS = append(hashS, hex.EncodeToString(it))
			}
			logger.Debug.Printf("File hash %s %s", strings.Join(hashS, " "), it.path)
			rs.result[it.path] = it.hash
		}
		resultCh <- rs
	}()
	wg.Wait()
	close(hashCh)

	rs := <-resultCh
	return rs.result, rs.errCount
}

func calcFileHash(path string) ([][]byte, error) {
	var rs [][]byte
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

	rs = append(rs, h1.Sum(make([]byte, 0, 16)))
	rs = append(rs, h2.Sum(make([]byte, 0, 16)))
	rs = append(rs, h3.Sum(make([]byte, 0, 16)))

	return rs, nil
}

// 扫描目录，深度优先
func _scanRecursive(path string, rsChan chan interface{}, m map[string]struct{}) {
	if _, ok := m[path]; ok {
		return
	}
	m[path] = struct{}{}
	entries, e := os.ReadDir(path)
	if e != nil {
		rsChan <- errors.Wrap(e, "ReadDir failed")
	} else {
		for _, entry := range entries {
			fileName := entry.Name()
			if strings.HasPrefix(fileName, ".") || fileHashScanBlackList[strings.ToLower(fileName)] {
				continue
			}
			if !entry.IsDir() {
				rsChan <- filepath.Join(path, fileName)
			} else {
				_scanRecursive(filepath.Join(path, fileName), rsChan, m)
			}
		}
	}
}
