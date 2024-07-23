package chunkupload

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
    "github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/infra/logctx"
)

// UploadDirectory will pack files in the directory to tar.gz stream and upload it
func UploadDirectory(ctx context.Context, dir string, fileFilter Filter, params Params, concurrentNumber int) error {
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

	var uw = NewUploadWriter(ctx, _ChunkSize, concurrentNumber, func(chunkId int, data []byte) error {
		return api.UploadCheckFiles(api.DefaultClient(), params.TaskId, params.SubtaskId, chunkId, bytes.NewReader(data))
	})
	e = dirPacker(ctx, dir, fileFilter, uw)
	if e != nil {
		return e
	}
	return uw.Close()
}
