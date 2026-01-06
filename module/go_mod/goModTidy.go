package go_mod

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/scanerr"
	"go.uber.org/zap"
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
	logger.Debug("go mod tidy :" + path)
	//记录次数
	var count int = 0
	var isourceError bool = false
again:

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stderr = &stdErr
	cmd.Dir = path
	if err := cmd.Start(); err != nil {
		logctx.Use(ctx).Error("Command finished with error" + err.Error())
		return err
	}
	if err := cmd.Wait(); err != nil {
		// 命令执行失败，先看 stderr 里是否包含需要特殊处理的错误（例如设置 GOPRIVATE）
		msg := stdErr.String()

		// 如果包含特定错误，尝试修复后重试一次
		if strings.Contains(msg, sourceError) {
			if e := goModTidyError(ctx, msg); e != nil {
				return e
			}
			isourceError = true
		}
		count++
		//如果因为设置GOPRIVATE失败，重试3次后还是失败，则返回错误
		if count == 3 {

			// 普通失败，记录基础扫描错误
			scanerr.Add(ctx, scanerr.Param{
				Kind:    "auto_build_error",
				Content: msg,
			})
			return errors.New(msg)
		} else if count == 2 && !isourceError {
			// 普通失败，记录基础扫描错误
			scanerr.Add(ctx, scanerr.Param{
				Kind:    "auto_build_error",
				Content: msg,
			})
			//如果因为不是因为设置GOPRIVATE失败，重试2次后还是失败，则返回错误，避免无限重试
			return errors.New(msg)
		}

		// 清空上一次的 stderr，再重试
		stdErr.Reset()
		goto again

	}

	// 走到这里说明 go mod tidy 成功了
	// 即使命令成功，stderr 也可能只是一些警告信息，这里只打日志，不触发基础扫描
	if stdErr.Len() > 0 {
		logger.Warn("go mod tidy stderr (ignored, command succeeded)",
			zap.String("stderr", stdErr.String()))
	}

	return nil
}
