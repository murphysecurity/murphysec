package chunkupload

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/klauspost/pgzip"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils"
	"golang.org/x/sync/errgroup"
)

type Params struct {
	TaskId    string
	SubtaskId string
}

const _ChunkSize = 4 * 1024 * 1024

func fileStreamer(ctx context.Context, path string, writer io.Writer) (e error) {
	var (
		logger = logctx.Use(ctx).Sugar()
		f      *os.File
	)
	f, e = os.Open(path)
	if e != nil {
		return e
	}
	defer func() { utils.LogCloseErr(logger, "file", f) }()
	info, e := os.Stat(path)
	if e != nil {
		return fmt.Errorf("fileStreamer: get fileinfo failed, %w", e)
	}
	logger.Infof("begin")
	defer func() { logger.Warnf("end with error: %v", e) }()
	gzipWriter := pgzip.NewWriter(writer)
	defer func() { utils.LogCloseErr(logger, "gzip", gzipWriter) }()
	tarWriter := tar.NewWriter(gzipWriter)
	defer func() { utils.LogCloseErr(logger, "tar", tarWriter) }()
	e = tarWriter.WriteHeader(&tar.Header{
		Name: filepath.Base(path),
		Mode: 666,
		Size: info.Size(),
	})
	if e != nil {
		return fmt.Errorf("fileStreamer: write header %w", e)
	}
	_, e = io.Copy(tarWriter, f)
	if e != nil {
		return e
	}
	return nil
}

func dirPacker(ctx context.Context, dir string, filter Filter, writer io.Writer) (_returnErr error) {
	var (
		e      error
		logger = logctx.Use(ctx).Sugar().Named("dirPacker")
	)
	logger.Infof("begin")
	defer func() { logger.Warnf("end with error: %v", _returnErr) }()

	gzipWriter := pgzip.NewWriter(writer)
	defer func() { utils.LogCloseErr(logger, "gzip", gzipWriter) }()
	tarWriter := tar.NewWriter(gzipWriter)
	defer func() { utils.LogCloseErr(logger, "tar", tarWriter) }()

	var putFile = func(p string, entry fs.DirEntry) error {
		var (
			e error
		)
		rel, e := filepath.Rel(dir, p)
		if e != nil {
			return fmt.Errorf("eval rel path: %w", e)
		}
		rel = filepath.ToSlash(rel)
		info, e := entry.Info()
		if e != nil {
			return fmt.Errorf("get info: %w", e)
		}
		f, e := os.Open(p)
		if e != nil {
			return fmt.Errorf("read file: %w", e)
		}
		defer func() { utils.LogCloseErr(logger, "file", f) }()

		var header tar.Header
		header.Name = rel
		header.Mode = 0666
		header.Size = info.Size()

		e = tarWriter.WriteHeader(&header)
		if e != nil {
			return fmt.Errorf("write header: %w", e)
		}

		_, e = io.Copy(tarWriter, f)
		if e != nil {
			return fmt.Errorf("write data: %w", e)
		}

		return nil
	}

	e = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		var (
			vote FilterVote
			e    error
		)
		vote, e = filter(path, d)
		if e != nil {
			return e
		}
		switch vote {
		case FilterSkip:
			return nil
		case FilterSkipDir:
			return filepath.SkipDir
		case FilterAdd:
		default:
			panic("bad value")
		}
		if d == nil {
			return fmt.Errorf("fs.DirEntry is nil")
		}
		if !d.Type().IsRegular() || d.Type().Type()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}
		if d.IsDir() {
			return nil
		} // ignore directory
		e = putFile(path, d)
		if e != nil {
			logger.Errorf("put file error: %s, %v", path, e)
			return e
		}
		return nil
	})
	if e != nil {
		return e
	}
	return e
}

func chunkUploadRoutine(ctx context.Context, params Params, reader io.Reader) error {
	goroutineNumber := min(runtime.NumCPU(), 1)
	var (
		// e          error
		eg, ec    = errgroup.WithContext(ctx)
		logger    = logctx.Use(ctx).Sugar().Named("chunkUploader")
		uploading = true
		bufferCh  = make(chan map[any]any, 1)
	)

	logger.Infof("begin")
	defer func() { logger.Infof("end") }()

	for range goroutineNumber {
		eg.Go(func() error {
			for bufInfo := range bufferCh {
				buffer := bufInfo["buffer"].(*bytes.Buffer)
				chunkId := bufInfo["chunkId"].(int)
				logger.Error("chunkId======" + strconv.Itoa(chunkId))
				err := api.UploadCheckFiles(api.DefaultClient(), params.TaskId, params.SubtaskId, chunkId, bytes.NewReader(buffer.Bytes()))
				if err != nil {
					logger.Error("api.UploadCheckFiles error:" + err.Error())
					return err
				}
			}
			return nil
		})
	}
	eg.Go(func() error {
		defer close(bufferCh)
		var chunkId = 0
		for uploading {
			var buf bytes.Buffer
			bufferInfo := make(map[any]any)
			chunkId++
			_, err := io.CopyN(&buf, reader, _ChunkSize)
			if err == io.EOF {
				uploading = false
				logger.Debugf("reader EOF")
			}
			bufferInfo["buffer"] = &buf
			bufferInfo["chunkId"] = chunkId
			select {
			case bufferCh <- bufferInfo:
			case <-ec.Done():
				break
			}
		}
		return nil
	})

	return eg.Wait()
	// for uploading {
	// 	chunkId++
	// 	e = ctx.Err()
	// 	if e != nil {
	// 		return e
	// 	}
	// 	_, e = io.CopyN(buf, reader, _ChunkSize)
	// 	if e == io.EOF {
	// 		e = nil
	// 		uploading = false
	// 		logger.Debugf("reader EOF")
	// 	}
	// 	if e != nil {
	// 		logger.Warnf("error during reading: %v", e)
	// 		return e
	// 	}
	// 	e = ctx.Err()
	// 	if e != nil {
	// 		return e
	// 	}
	// 	// uploading
	// 	e = api.UploadCheckFiles(api.DefaultClient(), params.TaskId, params.SubtaskId, chunkId, bytes.NewReader(buf.Bytes()))
	// 	if e != nil {
	// 		return e
	// 	}
	// 	buf.Reset()
	// }
	// return e
}

func checkDirValid(ctx context.Context, dir string) error {
	var (
		e      error
		info   os.FileInfo
		logger = logctx.Use(ctx).Sugar()
	)
	logger.Debugf("UploadDirectory: %s", dir)
	info, e = os.Stat(dir)
	if e != nil {
		logger.Warnf("stat: %v", e)
		return ErrDirInvalid
	}
	if !info.IsDir() {
		return ErrDirInvalid
	}
	return nil
}
