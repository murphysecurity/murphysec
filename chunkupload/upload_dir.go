package chunkupload

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path/filepath"

	ctxio "github.com/jbenet/go-context/io"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"golang.org/x/sync/errgroup"
)

// UploadDirectory will pack files in the directory to tar.gz stream and upload it
func UploadDirectory(ctx context.Context, dir string, fileFilter Filter, params Params) error {
	var (
		e      error
		logger = logctx.Use(ctx).Sugar()
	)
	logger.Infof("UploadDirectory, %s", dir)
	if fileFilter == nil {
		fileFilter = uploadAll
	}

	e = checkDirValid(ctx, dir)
	if e != nil {
		return e
	}
	dir, e = filepath.Abs(dir)
	if e != nil {
		return fmt.Errorf("eval abs path: %w", e)
	}

	eg, ec := errgroup.WithContext(ctx)
	pr, pw := io.Pipe()
	// create contextual io, avoid deadlock

	contextualReader := ctxio.NewReader(ec, pr)
	contextualWriter := ctxio.NewWriter(ec, pw)
	var writePipeBuffer = bufio.NewWriterSize(contextualWriter, 4*1024*1024)
	eg.Go(func() error { return chunkUploadRoutine(ctx, params, contextualReader) })
	eg.Go(func() (e error) {
		defer func() {
			if ee := writePipeBuffer.Flush(); ee != nil && e == nil {
				e = ee
			}
			if ee := pw.Close(); ee != nil && e == nil {
				e = pw.Close()
			}
		}()
		return dirPacker(ctx, dir, fileFilter, writePipeBuffer)
	})

	e = eg.Wait()
	if e != nil {
		logger.Warnf("UploadDirectory failed, %v", e)
		return e
	}
	return nil
}
