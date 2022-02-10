package inspector

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
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
		hashStr[hex.EncodeToString(s)] = struct{}{}
	}
	if ctx.TaskSource == api.TaskSourceCli {
		fmt.Println(fmt.Sprintf("共成功扫描 %d 个文件", errCount))
	}
	// api request object
	req := ctx.getApiRequestObj()
	for s := range hashStr {
		req.FileHashList = append(req.FileHashList, api.FileHash{Hash: s})
	}
	return api.SendDetectHash(req)
}

const ScanConcurrency = 2

// todo: refactor

func FileHashInspectScan(projectDir string) (map[string][]byte, int) {

	fch := make(chan interface{}, 100)
	go func() {
		_scanRecursive(projectDir, fch, map[string]struct{}{})
		close(fch)
	}()
	wg := sync.WaitGroup{}
	type hashResult struct {
		path string
		hash []byte
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
		result   map[string][]byte
	}
	resultCh := make(chan result, 1)
	go func() {
		rs := result{
			errCount: 0,
			result:   map[string][]byte{},
		}
		for it := range hashCh {
			if it.err != nil {
				rs.errCount++
				logger.Warn.Println("Calc hash failed: %v", it.err)
				continue
			}
			logger.Debug.Printf("File hash %s %s", hex.EncodeToString(it.hash), it.path)
			rs.result[it.path] = it.hash
		}
		resultCh <- rs
	}()
	wg.Wait()
	close(hashCh)

	rs := <-resultCh
	return rs.result, rs.errCount
}

func calcFileHash(path string) ([]byte, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, errors.Wrap(e, fmt.Sprintf("Open file failed when calc file hash: %s", path))
	}
	defer f.Close()
	m := md5.New()
	if _, e := io.Copy(m, f); e != nil {
		return nil, errors.Wrap(e, fmt.Sprintf("Calc file hash failed: %s", path))
	}
	return m.Sum(make([]byte, 0, 16)), nil
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
