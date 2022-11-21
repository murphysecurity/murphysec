package chunkupload

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Params struct {
	TaskId    string
	SubtaskId string
}

type FilterVote int

const (
	_ FilterVote = iota
	FilterAdd
	FilterSkip
	FilterSkipDir
)

type Filter func(path string, entry fs.DirEntry) (FilterVote, error)

// UploadDirectory will pack files in the directory to tar.gz stream and upload it
func UploadDirectory(ctx context.Context, dir string, fileFilter Filter, params Params) error {
	var (
		e      error
		logger = logctx.Use(ctx).Sugar()
	)
	e = checkDirValid(ctx, dir)
	if e != nil {
		return e
	}
	dir, e = filepath.Abs(dir)
	if e != nil {
		logger.Errorf("%v", e)
		return ErrEvalAbsPath
	}
	r, w := io.Pipe()
	defer w.Close()
	gzipStream := gzip.NewWriter(w)
	tarStream := tar.NewWriter(gzipStream)
	eg, ec := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer func() {
			var e error
			e = tarStream.Close()
			if e != nil {
				logger.Warnf("tar stream close with error: %v", e)
			}
			e = gzipStream.Close()
			if e != nil {
				logger.Warnf("gzip stream close with error: %v", e)
			}
		}()
		return fileWalker(ec, dir, fileFilter, tarStream)
	})
	eg.Go(func() error {
		e := chunkUploadStreamer(ec, r, params)
		if e != nil {

		}
	})

}

func fileWalker(ctx context.Context, dir string, filter Filter, writer *tar.Writer) error {
	var (
		vote   FilterVote
		e      error
		logger = logctx.Use(ctx).Sugar()
	)
	logger.Infof("fileWalker: begin")
	defer func() { logger.Infof("fileWalker: end, %v", e) }()

	var putFile = func(p string, entry fs.DirEntry) error {
		var (
			rel    string
			e      error
			header tar.Header
			info   fs.FileInfo
			f      *os.File
		)
		rel, e = filepath.Rel(dir, p)
		if e != nil {
			return e
		}
		rel = filepath.ToSlash(rel)
		info, e = entry.Info()
		if e != nil {
			return e
		}
		header.Name = rel
		header.Mode = 0666
		header.Size = info.Size()
		e = writer.WriteHeader(&header)
		if e != nil {
			return e
		}
		f, e = os.Open(p)
		if e != nil {
			return e
		}
		defer func() {
			if e := f.Close(); e != nil {
				logger.Warnf("close file failed, %v", e)
			}
		}()
		return nil
	}
	e = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if e == nil {
			return e
		}
		if d == nil {
			return errors.New("fs entry is nil")
		}
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
			logger.Debugf("put file: %s", path)
			e = putFile(path, d)
			if e != nil {
				return e
			}
		default:
			panic("bad")
		}
		return nil
	})
	return e
}

const ChunkSize = 16 * 1024 * 1024

func chunkUploadStreamer(ctx context.Context, reader io.Reader, params Params) (e error) {
	var (
		buf       = &bytes.Buffer{}
		logger    = logctx.Use(ctx).Sugar()
		n64       int64
		chunkId   int
		uploading = true
	)

	logger.Infof("chunkUploadStreamer: begin")
	defer func() { logger.Infof("chunkUploadStreamer: end, %v", e) }()

	for uploading {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		buf.Reset()
		chunkId++
		n64, e = io.CopyN(buf, reader, ChunkSize)
		logger.Debugf("chunk copy[%d] %d bytes, %v", chunkId, n64, e)
		if e != nil {
			if e == io.EOF {
				uploading = false
			} else {
				return e
			}
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		e = api.UploadCheckFiles(api.DefaultClient(), params.TaskId, params.SubtaskId, chunkId, bytes.NewReader(buf.Bytes()))
		if e != nil {
			return e
		}
	}
	logger.Infof("chunk copy ended")
	return nil
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

}

//go:generate stringer -type uploadErr -linecomment
type uploadErr int

const (
	_              uploadErr = iota
	ErrDirInvalid            // uploader: dir invalid
	ErrEvalAbsPath           // uploader: cannot evaluate absolute path
)

func (i uploadErr) Error() string {
	return i.String()
}
