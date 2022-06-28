package inspector

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	ctxio "github.com/jbenet/go-context/io"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func BinScan(ctx context.Context) error {
	var e error
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.Display(display.MsgInfo, fmt.Sprint("项目名称：", scanTask.ProjectName))
	if e := createTaskC(ctx); e != nil {
		return e
	}
	if e := binScanUploadFile(ctx); e != nil {
		return e
	}

	// 开始扫描
	if e := api.StartCheckTaskType(scanTask.TaskId, scanTask.Kind); e != nil {
		logger.Err.Println("StartCheck failed.", e.Error())
		ui.Display(display.MsgError, fmt.Sprint("开始扫描失败 ", e.Error()))
		return e
	}
	// 等待返回结果
	var r *model.TaskScanResponse
	ui.WithStatus(display.StatusRunning, "已提交，正在扫描...", func() {
		r, e = api.QueryResult(scanTask.TaskId)
	})
	if e != nil {
		logger.Err.Println("QueryResult failed.", e.Error())
		fmt.Println("扫描失败", e.Error())
		ui.Display(display.MsgError, fmt.Sprint("扫描失败，", e.Error()))
		return e
	} else {
		ui.Display(display.MsgNotice, fmt.Sprintf("项目扫描完成，依赖数：%d，漏洞数：%d\n", r.DependenciesCount, r.IssuesCompsCount))
	}
	return nil
}

func binScanUploadFile(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.UpdateStatus(display.StatusRunning, "正在上传文件...")
	defer ui.ClearStatus()

	pathCh := make(chan string, 10)
	g, goCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return scanBinaryFile(goCtx, scanTask.ProjectDir, pathCh) })
	r, w := io.Pipe()
	g.Go(func() error { return packFileToTgzStream(goCtx, pathCh, scanTask.ProjectDir, w) })
	g.Go(func() error { return uploadTgzChunk(goCtx, r) })
	if e := g.Wait(); e != nil {
		logger.Err.Println(e)
		ui.Display(display.MsgError, fmt.Sprint("文件上传失败：", e.Error()))
		return e
	}
	ui.Display(display.MsgInfo, "文件上传成功")
	return nil
}

func uploadTgzChunk(goctx context.Context, r io.Reader) error {
	task := model.UseScanTask(goctx)
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
		e = api.UploadChunk(task.TaskId, counter, io.LimitReader(bytes.NewReader(buf), int64(n)))
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
		rp := filepath.ToSlash(must.A(filepath.Rel(baseDir, s)))
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
