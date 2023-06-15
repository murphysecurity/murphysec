package chunkupload

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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
	gzipWriter := gzip.NewWriter(writer)
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

	gzipWriter := gzip.NewWriter(writer)
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
		if d == nil {
			return fmt.Errorf("fs.DirEntry is nil")
		}
		if !d.Type().IsRegular() || d.Type().Type()&os.ModeSymlink == os.ModeSymlink {
			return nil
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
	var (
		e         error
		logger    = logctx.Use(ctx).Sugar().Named("chunkUploader")
		buf       = &bytes.Buffer{}
		uploading = true
		chunkId   int
	)
	logger.Infof("begin")
	defer func() { logger.Infof("end, %v", e) }()

	for uploading {
		chunkId++
		e = ctx.Err()
		if e != nil {
			return e
		}
		_, e = io.CopyN(buf, reader, _ChunkSize)
		if e == io.EOF {
			e = nil
			uploading = false
			logger.Debugf("reader EOF")
		}
		if e != nil {
			logger.Warnf("error during reading: %v", e)
			return e
		}
		e = ctx.Err()
		if e != nil {
			return e
		}
		// uploading
		e = api.UploadCheckFiles(api.DefaultClient(), params.TaskId, params.SubtaskId, chunkId, bytes.NewReader(buf.Bytes()))
		if e != nil {
			return e
		}
		buf.Reset()
	}
	return e
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
