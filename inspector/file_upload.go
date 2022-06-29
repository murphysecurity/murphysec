package inspector

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func UploadCodeFile(ctx context.Context) error {
	task := model.UseScanTask(ctx)
	codeFiles := scanCodeFile(ctx)
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
				logger := Logger.With(zap.String("path", p))
				if e != nil {
					logger.Warn("Stat failed", zap.Error(e))
					continue
				}
				rp, e := filepath.Rel(task.ProjectDir, p)
				if e != nil {
					logger.Warn("Get relative-path failed", zap.Error(e))
					continue
				}
				rp = filepath.ToSlash(rp)
				if strings.HasPrefix("./", rp) {
					logger.Warn("Bad prefix")
					continue
				}
				f, e := os.Open(p)
				if e != nil {
					logger.Warn("Open file failed", zap.Error(e))
					continue
				}
				e = tarWriter.WriteHeader(&tar.Header{
					Name: rp,
					Size: info.Size(),
					Mode: 0666,
				})
				if e != nil {
					logger.Warn("Write header failed", zap.Error(e))
					return e
				}
				_, e = io.Copy(tarWriter, f)
				if e != nil {
					logger.Warn("Write file failed", zap.Error(e))
					return e
				}
				utils.CloseLogErrZap(f, logger)
			}
			return nil
		}()
		if e != nil {
			failure = true
			Logger.Error("Pkg files failed", zap.Error(e))
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
				Logger.Info("read zero byte")
				continue
			}
			e = api.UploadChunk(task.TaskId, counter, io.LimitReader(bytes.NewReader(buf), int64(n)))
			if e != nil {
				Logger.Error("Upload file failed", zap.Error(e))
				failure = true
				return
			} else {
				Logger.Info("block sent")
			}
		}
	}()
	wg.Wait()
	if failure {
		return errors.New("Upload file failed")
	}
	return nil
}

func scanCodeFile(ctx context.Context) []string {
	task := model.UseScanTask(ctx)
	Logger.Debug("Start scan code files", zap.String("project_dir", task.ProjectDir))
	fileSet := map[string]struct{}{}
	e := filepath.WalkDir(task.ProjectDir, func(path string, d fs.DirEntry, err error) error {
		if d == nil || err != nil {
			Logger.Warn("Walk error", zap.Error(err), zap.String("path", path))
			return nil
		}
		if dirShouldIgnore(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}
		info, e := d.Info()
		if e != nil {
			Logger.Warn("Get file info failed", zap.Error(e), zap.String("path", path))
			return nil
		}
		if info.Size() < 32 || info.Size() > 4*1024*1024 {
			return nil
		}
		if codeFileShouldUpload(d.Name()) {
			fileSet[path] = struct{}{}
		}
		return nil
	})
	Logger.Info("Code file scan finished", zap.Int("total", len(fileSet)))
	if e != nil {
		Logger.Error("Error happened during code file scan", zap.Error(e))
	}
	var rs []string
	for s := range fileSet {
		rs = append(rs, s)
	}
	return rs
}

func codeFileShouldUpload(p string) bool {
	ext := filepath.Ext(p)
	_, ok := __uploadFileExt[ext]
	return ok
}

var __uploadFileExt = map[string]struct{}{}

//go:embed file_upload_ext
var __uploadFileExtList string

func init() {
	for _, s := range strings.Split(__uploadFileExtList, "\n") {
		s = strings.TrimSpace(s)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		__uploadFileExt[s] = struct{}{}
	}
}
