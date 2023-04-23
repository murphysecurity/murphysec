package chunkupload

import (
	"context"
	"fmt"
	ctxio "github.com/jbenet/go-context/io"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"golang.org/x/sync/errgroup"
	"io"
	"path/filepath"
)

func UploadFile(ctx context.Context, path string, params Params) error {
	var (
		e      error
		logger = logctx.Use(ctx).Sugar()
	)
	logger.Infof("UploadFile, %s", path)
	path, e = filepath.Abs(path)
	if e != nil {
		return fmt.Errorf("eval abs path: %w", e)
	}

	eg, ec := errgroup.WithContext(ctx)
	pr, pw := io.Pipe()
	// create contextual io, avoid deadlock
	contextualReader := ctxio.NewReader(ec, pr)
	contextualWriter := ctxio.NewWriter(ec, pw)

	eg.Go(func() error { return chunkUploadRoutine(ctx, params, contextualReader) })
	eg.Go(func() error {
		defer func() { _ = pw.Close() }()
		return fileStreamer(ctx, path, contextualWriter)
	})

	e = eg.Wait()
	if e != nil {
		logger.Warnf("UploadFile failed, %v", e)
		return e
	}

	return nil
}
