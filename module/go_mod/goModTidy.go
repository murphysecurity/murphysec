package go_mod

import (
	"bytes"
	"context"
	"errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
)

const (
	sourceError = "no secure protocol found for repository"
)

func setPrivate(ctx context.Context, privateUrl string) error {

	privateUrlList := os.Getenv("GOPRIVATE")
	privateUrlList = privateUrlList + "," + privateUrl
	if err := os.Setenv("GOPRIVATE", privateUrlList); err != nil {
		logctx.Use(ctx).Error("Failed to set GOPRIVATE", zap.Error(err))
		return err
	}
	logctx.Use(ctx).Info("set GOPRIVATE = " + privateUrlList)
	return nil
}
func goModTidyError(ctx context.Context, msg string) error {
	if strings.Contains(msg, sourceError) {
		var privateUrl string
		if u := strings.Split(msg, "/"); len(u) > 1 {
			privateUrl = u[0]
			if privateUrl != "" {
				return setPrivate(ctx, privateUrl)
			}
		}
	}
	return nil
}
func goModTidy(ctx context.Context, path string) error {
	logger := logctx.Use(ctx)
	var stdErr bytes.Buffer
	var againBol = false
	logger.Debug("go mod tidy :" + path)
again:

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stderr = &stdErr
	cmd.Dir = path
	if err := cmd.Start(); err != nil {
		logctx.Use(ctx).Error("Command finished with error" + err.Error())
		return err
	}
	cmd.Wait()

	if againBol {
		return errors.New(stdErr.String())
	}
	if stdErr.Len() > 0 {
		if err := goModTidyError(ctx, stdErr.String()); err != nil {
			return err
		} else {
			againBol = true
			goto again
		}
	}
	return nil
}
