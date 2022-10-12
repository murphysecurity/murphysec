package inspector

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	ctxio "github.com/jbenet/go-context/io"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/murphysecurity/murphysec/view"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func BinScan(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	view.ProjectName(ui, scanTask.ProjectName)
	if e := createTaskC(ctx); e != nil {
		return e
	}
	if e := binScanUploadFile(ctx); e != nil {
		return e
	}

	// 开始扫描
	if e := startCheckC(ctx); e != nil {
		return e
	}
	// 等待返回结果
	if e := queryResultC(ctx); e != nil {
		return e
	}
	r := scanTask.ScanResult
	view.DisplayScanResultSummary(ui, r.DependenciesCount, r.IssuesCompsCount)
	view.DisplayScanResultReport(ui, scanTask.ScanResult.ReportURL())

	return nil
}

func binScanUploadFile(ctx context.Context) error {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	defer view.FileUploading(ui)()
	pathCh := make(chan string, 10)
	g, goCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return scanBinaryFile(goCtx, scanTask.ProjectDir, pathCh) })
	r, w := io.Pipe()
	g.Go(func() error { return packFileToTgzStream(goCtx, pathCh, scanTask.ProjectDir, w) })
	g.Go(func() error { return uploadTgzChunk(goCtx, r) })
	if e := g.Wait(); e != nil {
		Logger.Error("Upload failed", zap.Error(e))
		view.FileUploadFailed(ui, e)
		return e
	}
	view.FileUploadSucceeded(ui)
	return nil
}

func uploadTgzChunk(goctx context.Context, r io.Reader) error {
	logger := Logger.Named("tgz-uploader")
	task := model.UseScanTask(goctx)
	counter := 0
	for {
		if goctx.Err() != nil {
			logger.Info("context error:", zap.Error(goctx.Err()))
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
		logger.Debug("write", zap.Int("length", n))
		e = api.UploadChunk(task.TaskId, counter, io.LimitReader(bytes.NewReader(buf), int64(n)))
		if e != nil {
			return errors.Wrap(e, "文件上传失败")
		} else {
			logger.Info("block sent")
		}
	}
	return nil
}

func packFileToTgzStream(ctx context.Context, fileNameCh chan string, baseDir string, w io.WriteCloser) error {
	logger := Logger.Named("file-tgz-pack")
	defer logger.Info("terminated")
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
			logger.Warn("Stat file failed", zap.String("file", s), zap.Error(e))
			continue
		}
		rp := filepath.ToSlash(must.A(filepath.Rel(baseDir, s)))
		if rp == "." {
			rp = filepath.Base(baseDir)
		}
		f, e := os.Open(s)
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
			utils.CloseLogErrZap(f, logger)
			return e
		}
		_, e = io.Copy(tarWriter, f)
		utils.CloseLogErrZap(f, logger)
		if e != nil {
			return e
		}
	}
	return nil
}

var _ErrScanBinaryWalkStop = errors.New("_ErrScanBinaryWalkStop")

func scanBinaryFile(ctx context.Context, dir string, pathCh chan string) error {
	var logger = Logger.Named("binary-file-walker")
	logger.Info("", zap.String("dir", dir))
	defer logger.Info("terminated")
	defer close(pathCh)
	info, e := os.Stat(dir)
	if e != nil {
		logger.Error("Stat dir failed", zap.Error(e))
		return errors.Wrap(e, "ScanBinaryFile")
	}
	if !info.IsDir() {
		logger.Warn("Not a dir")
		pathCh <- dir
		return nil
	}
	e = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info == nil || err != nil {
			logger.Error("walk error", zap.Error(e))
			return nil
		}
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || binaryScanDirBlackList[info.Name()]) {
			return filepath.SkipDir
		}
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") || !info.Mode().IsRegular() {
			return nil
		}
		select {
		case pathCh <- path:
		case <-ctx.Done():
			return _ErrScanBinaryWalkStop
		}
		return nil
	})
	if e == nil {
		return nil
	}
	if e == _ErrScanBinaryWalkStop {
		logger.Warn("Walk cancel")
		return nil
	}
	logger.Error("Walk error", zap.Error(e))
	return e
}

var binaryScanDirBlackList = map[string]bool{
	"node_modules": true,
}
