package inspector

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func UploadCodeFile(ctx *ScanContext) error {
	codeFiles := ScanCodeFile(ctx)
	if len(codeFiles) == 0 {
		return nil
	}
	r, pw := io.Pipe()
	w := bufio.NewWriterSize(pw, 4*1024*1024)
	gzipWriter := gzip.NewWriter(w)
	tarWriter := tar.NewWriter(gzipWriter)
	failure := false
	go func() {
		e := func() error {
			for _, p := range codeFiles {
				info, e := os.Stat(p)
				if e != nil {
					logger.Warn.Println("os.Stat file failed.", e.Error(), p)
					continue
				}
				rp, e := filepath.Rel(ctx.ProjectDir, p)
				if e != nil {
					logger.Warn.Println("get relative-path failed.", e.Error(), p)
					continue
				}
				rp = filepath.ToSlash(rp)
				if strings.HasPrefix("./", rp) {
					logger.Warn.Println("bad prefix of path", rp)
					continue
				}
				logger.Debug.Println("tar append relative-path:", rp)
				f, e := os.Open(p)
				if e != nil {
					logger.Warn.Println("Open file failed.", e.Error(), p)
					continue
				}
				e = tarWriter.WriteHeader(&tar.Header{
					Name: rp,
					Size: info.Size(),
					Mode: 0666,
				})
				if e != nil {
					logger.Warn.Println("Append tar header failed.", e.Error())
					return e
				}
				n, e := io.Copy(tarWriter, f)
				if e != nil {
					logger.Warn.Println("Append tar content failed.", e.Error())
					return e
				}
				logger.Debug.Println("Append tar", n, "bytes")
				e = f.Close()
				if e != nil {
					logger.Warn.Println("Close file failed.", e.Error())
				}
			}
			return nil
		}()
		if e != nil {
			failure = true
			logger.Err.Println("Pkg files failed.", e.Error())
		}
		must.Must(tarWriter.Close())
		must.Must(gzipWriter.Close())
		must.Must(w.Flush())
		must.Must(pw.Close())
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		counter := 0
		for {
			counter++
			buf := make([]byte, 4*1024*1024)
			n, e := r.Read(buf)
			if e != nil {
				if e == io.EOF && n == 0 {
					break
				}
				if e != io.EOF {
					return
				}
			}
			if n == 0 {
				logger.Info.Println("read zero byte")
				continue
			}
			e = api.UploadChunk(ctx.TaskId, counter, io.LimitReader(bytes.NewReader(buf), int64(n)))
			if e != nil {
				logger.Err.Println("Upload file failed.", e.Error())
				failure = true
				return
			} else {
				logger.Info.Println("block sent")
			}
		}
	}()
	wg.Wait()
	if failure {
		return errors.New("Upload file failed")
	}
	return nil
}

func ScanCodeFile(ctx *ScanContext) []string {
	logger.Debug.Println("Start scan code files:", ctx.ProjectDir, "...")
	fileSet := map[string]struct{}{}
	e := filepath.Walk(ctx.ProjectDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || folderNameBlackList[info.Name()]) {
			return filepath.SkipDir
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if info.Size() < 32 || info.Size() > 4*1024*1024 {
			return nil
		}
		if uploadFileExt[strings.TrimPrefix(filepath.Ext(info.Name()), ".")] {
			fileSet[path] = struct{}{}
		}
		return nil
	})
	logger.Debug.Println("Code file scan finished, total", len(fileSet))
	if e != nil {
		logger.Warn.Println("Error happened during code file scan,", e.Error())
	}
	var rs []string
	for s := range fileSet {
		rs = append(rs, s)
	}
	return rs
}

var uploadFileExt = map[string]bool{
	"c":    true,
	"C":    true,
	"cpp":  true,
	"h":    true,
	"hpp":  true,
	"cxx":  true,
	"c++":  true,
	"java": true,
	"xml":  true,
}

var folderNameBlackList = map[string]bool{
	"node_modules": true,
	"__MACOSX":     true,
}
