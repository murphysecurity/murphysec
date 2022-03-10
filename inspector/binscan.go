package inspector

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"strings"
)

func BinScan(projectDir string) error {
	{
		_, e := os.Stat(projectDir)
		if e != nil {
			fmt.Println("路径不存在或无效")
			return errors.New("Path doesn't exists")
		}
	}
	ctx := createBinaryTaskContext(projectDir)
	ctx.ProjectName = filepath.Base(projectDir)
	if e := createTask(ctx); e != nil {
		logger.Err.Println("create task failed.", e.Error())
		fmt.Println("创建任务失败！", e.Error())
		return e
	}
	fmt.Println("项目创建成功！", ctx.TaskId)
	rsCh := make(chan string, 10)
	go scanBinaryFile(ctx.ProjectDir, rsCh)
	r, w := io.Pipe()
	errCh := make(chan error, 10)
	go packFileToTgzStream(rsCh, projectDir, w, errCh)
	failure := false
	func() {
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
			logger.Debug.Println("write", n, "bytes")
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
	if failure {
		fmt.Println("文件上传失败")
		return errors.New("File upload failed.")
	}
	if e := api.StartCheck(ctx.TaskId, api.TaskKindBinary); e != nil {
		logger.Err.Println("StartCheck failed.", e.Error())
		fmt.Println("开始扫描失败")
	}
	fmt.Println("已提交，正在扫描...")
	if r, e := api.QueryResult(ctx.TaskId); e != nil {
		logger.Err.Println("QueryResult failed.", e.Error())
		fmt.Println("扫描失败", e.Error())
	} else {
		fmt.Printf("扫描成功！共计%d个组件，其中%d个存在风险\n", r.DependenciesCount, r.IssuesCompsCount)
	}
	return nil
}

func packFileToTgzStream(fch chan string, baseDir string, w io.WriteCloser, errCh chan error) {
	defer w.Close()
	bw := bufio.NewWriterSize(w, 1024*1024*4)
	defer bw.Flush()
	gzipWriter := gzip.NewWriter(bw)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	for s := range fch {
		info, e := os.Stat(s)
		if e != nil {
			logger.Err.Println("Stat file failed.", e.Error(), s)
			continue
		}
		rp := filepath.ToSlash(must.String(filepath.Rel(baseDir, s)))
		if rp == "." {
			rp = filepath.Base(baseDir)
		}
		f, e := os.Open(s)
		if e != nil {
			logger.Err.Println("Open file failed.", e.Error(), s)
			continue
		}
		e = tarWriter.WriteHeader(&tar.Header{
			Name: rp,
			Size: info.Size(),
			Mode: 0666,
		})
		if e != nil {
			logger.Err.Println(e.Error())
			errCh <- e
			return
		}
		n, e := io.Copy(tarWriter, f)
		f.Close()
		if e != nil {
			logger.Err.Println(e.Error())
			errCh <- e
			return
		}
		logger.Debug.Println("File size:", n)
	}
}

func scanBinaryFile(dir string, pathCh chan string) {
	defer close(pathCh)
	info, e := os.Stat(dir)
	if e != nil {
		logger.Err.Println("获取路径类型失败")
		return
	}
	if !info.IsDir() {
		pathCh <- dir
	}
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") || !info.Mode().IsRegular() {
			return nil
		}
		if info.Size() > 1*1024*1024*1024 {
			return nil
		}
		logger.Debug.Println("scan file:", path)
		pathCh <- path
		return nil
	})
}
