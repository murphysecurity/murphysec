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
	"murphysec-cli-simple/display"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/eg"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"strings"
)

func BinScan(ctx *ScanContext) error {
	ui := ctx.UI()
	e := func() error {
		ui.UpdateStatus(display.StatusRunning, "正在创建任务")
		defer ui.ClearStatus()
		return createTask(ctx)
	}()
	if e != nil {
		logger.Err.Println("create task failed.", e.Error())
		ui.Display(display.MsgError, fmt.Sprint("创建任务失败！"))
		if errors.Is(api.ErrTokenInvalid, e) {
			ui.Display(display.MsgError, "原因：当前 Token 无效")
		}
		return e
	}
	ui.Display(display.MsgInfo, fmt.Sprint("项目创建成功！", ctx.TaskId))
	e = func() error {
		ui.UpdateStatus(display.StatusRunning, "正在上传文件...")
		defer ui.ClearStatus()
		g := new(eg.EG)
		pathCh := make(chan string, 10)
		g.Go(func() { scanBinaryFile(ctx.ProjectDir, pathCh) })
		r, w := io.Pipe()
		g.Go(func() { packFileToTgzStream(pathCh, ctx.ProjectDir, w) })
		g.Go(func() { uploadTgzChunk(r, ctx) })
		return g.Wait()
	}()
	if e != nil {
		logger.Err.Println(e)
		ui.Display(display.MsgError, fmt.Sprint("文件上传失败：", e.Error()))
		return e
	}
	ui.Display(display.MsgInfo, "文件上传成功")
	if e := api.StartCheckTaskType(ctx.TaskId, ctx.Kind); e != nil {
		logger.Err.Println("StartCheck failed.", e.Error())
		ui.Display(display.MsgError, fmt.Sprint("开始扫描失败 ", e.Error()))
		return e
	}
	r, e := func() (*api.TaskScanResponse, error) {
		ui.UpdateStatus(display.StatusRunning, "已提交，正在扫描...")
		defer ui.ClearStatus()
		return api.QueryResult(ctx.TaskId)
	}()
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

func uploadTgzChunk(r io.Reader, ctx *ScanContext) {
	counter := 0
	for {
		counter++
		buf := make([]byte, 4*1024*1024)
		n, e := r.Read(buf)
		if e == io.EOF {
			break
		}
		must.Must(e)
		if n == 0 {
			continue
		}
		logger.Debug.Println("write", n, "bytes")
		e = api.UploadChunk(ctx.TaskId, counter, io.LimitReader(bytes.NewReader(buf), int64(n)))
		if e != nil {
			panic(errors.Wrap(e, "文件上传失败"))
		} else {
			logger.Info.Println("block sent", n)
		}
	}

}

func packFileToTgzStream(fileNameCh chan string, baseDir string, w io.WriteCloser) {
	defer w.Close()
	bw := bufio.NewWriterSize(w, 1024*1024*4)
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
		must.Must(e)
		_, e = io.Copy(tarWriter, f)
		must.Must(e)
	}
}

func scanBinaryFile(dir string, pathCh chan string) {
	defer close(pathCh)
	info, e := os.Stat(dir)
	must.Must(e)
	if !info.IsDir() {
		pathCh <- dir
	}
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			logger.Warn.Println("WalkErr:", err)
		}
		if info == nil {
			return nil
		}
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") || !info.Mode().IsRegular() {
			return nil
		}
		if info.Size() > 1*1024*1024*1024 {
			return nil
		}
		pathCh <- path
		return nil
	})
}
