package inspector

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	ctxio "github.com/jbenet/go-context/io"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"io/fs"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/display"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"strings"
)

func BinScan(ctx *ScanContext) error {
	ui := ctx.UI()
	var e error
	// 创建项目
	ui.WithStatus(display.StatusRunning, "正在创建任务", func() {
		e = createTask(ctx)
	})
	if e != nil {
		logger.Err.Println("create task failed.", e.Error())
		ui.Display(display.MsgError, fmt.Sprint("创建任务失败！"))
		if errors.Is(api.ErrTokenInvalid, e) {
			ui.Display(display.MsgError, "原因：当前 Token 无效")
		}
		return e
	}
	ui.Display(display.MsgInfo, fmt.Sprint("项目创建成功！", ctx.TaskId))

	// 上传文件
	ui.WithStatus(display.StatusRunning, "正在上传文件...", func() {
		pathCh := make(chan string, 10)
		g, goCtx := errgroup.WithContext(context.Background())
		g.Go(func() error { return scanBinaryFile(goCtx, ctx.ProjectDir, pathCh) })
		r, w := io.Pipe()
		g.Go(func() error { return packFileToTgzStream(goCtx, pathCh, ctx.ProjectDir, w) })
		g.Go(func() error { return uploadTgzChunk(goCtx, r, ctx) })
		e = g.Wait()
	})
	if e != nil {
		logger.Err.Println(e)
		ui.Display(display.MsgError, fmt.Sprint("文件上传失败：", e.Error()))
		return e
	}
	ui.Display(display.MsgInfo, "文件上传成功")
	// 开始扫描
	if e := api.StartCheckTaskType(ctx.TaskId, ctx.Kind); e != nil {
		logger.Err.Println("StartCheck failed.", e.Error())
		ui.Display(display.MsgError, fmt.Sprint("开始扫描失败 ", e.Error()))
		return e
	}
	// 等待返回结果
	var r *api.TaskScanResponse
	ui.WithStatus(display.StatusRunning, "已提交，正在扫描...", func() {
		r, e = api.QueryResult(ctx.TaskId)
	})
	if e != nil {
		logger.Err.Println("QueryResult failed.", e.Error())
		fmt.Println("扫描失败", e.Error())
		ui.Display(display.MsgError, fmt.Sprint("扫描失败，", e.Error()))
		return e
	} else {
		ui.Display(display.MsgNotice, fmt.Sprintf("项目扫描成功，依赖数：%d，漏洞数：%d\n", r.DependenciesCount, r.IssuesCompsCount))
	}
	return nil
}

func uploadTgzChunk(goctx context.Context, r io.Reader, ctx *ScanContext) error {
	counter := 0
	for {
		if goctx.Err() != nil {
			logger.Info.Println("context error:", goctx.Err())
			break
		}
		counter++
		buf := make([]byte, 4*1024*1024)
		n, e := ctxio.NewReader(goctx, r).Read(buf)
		if e == io.EOF || e == context.Canceled {
			break
		}
		if e != nil {
			return e
		}
		must.True(n > 0)
		logger.Debug.Println("write", n, "bytes")
		e = api.UploadChunk(ctx.TaskId, counter, io.LimitReader(bytes.NewReader(buf), int64(n)))
		if e != nil {
			return errors.Wrap(e, "文件上传失败")
		} else {
			logger.Info.Println("block sent", n)
		}
	}
	return nil
}

func packFileToTgzStream(ctx context.Context, fileNameCh chan string, baseDir string, w io.WriteCloser) error {
	defer w.Close()
	bw := bufio.NewWriterSize(ctxio.NewWriter(ctx, w), 1024*1024*4)
	defer bw.Flush()
	gzipWriter := gzip.NewWriter(bw)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	for s := range fileNameCh {
		info, e := os.Stat(s)
		if e != nil {
			logger.Warn.Println("Stat file failed.", e.Error(), s)
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
			return e
		}
		_, e = io.Copy(tarWriter, f)
		if e != nil {
			return e
		}
	}
	return nil
}

var _ErrScanBinaryWalkStop = errors.New("_ErrScanBinaryWalkStop")

func scanBinaryFile(ctx context.Context, dir string, pathCh chan string) error {
	defer close(pathCh)
	info, e := os.Stat(dir)
	if e != nil {
		return errors.Wrap(e, "ScanBinaryFile")
	}
	if !info.IsDir() {
		pathCh <- dir
	}
	e = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info == nil || err != nil {
			logger.Info.Println("Error during filepath walk:", err)
			return nil
		}
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || binaryScanDirBlackList[info.Name()]) {
			return filepath.SkipDir
		}
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") || !info.Mode().IsRegular() {
			return nil
		}
		if info.Size() > 1*1024*1024*1024 {
			return nil
		}
		select {
		case pathCh <- path:
		case <-ctx.Done():
			return _ErrScanBinaryWalkStop
		}
		return nil
	})
	if e == _ErrScanBinaryWalkStop {
		logger.Info.Println("ScanBinaryFile filepath walk cancel")
		return nil
	}
	if e != nil {
		logger.Warn.Println("Error happens during filepath walk:", e.Error())
		return errors.Wrap(e, "ScanBinaryFile")
	}
	return nil
}

var binaryScanDirBlackList = map[string]bool{
	"node_modules": true,
}
