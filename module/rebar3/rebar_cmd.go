package rebar3

import (
	"bytes"
	"context"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/logpipe"
	"io"
	"os/exec"
)

var ErrCallRebar3Command = errors.New("Execute rebar3 command failed")

var __rebarVersionCached = ""
var __rebarVersionErrorCached error

func GetRebar3Version(ctx context.Context) (string, error) {
	if __rebarVersionErrorCached != nil {
		return "", __rebarVersionErrorCached
	}
	if __rebarVersionCached != "" {
		return __rebarVersionCached, nil
	}
	var logger = logctx.Use(ctx)
	var cmd = exec.Command("rebar3", "version")
	logger.Sugar().Infof("Execute command: %s", cmd.String())
	data, e := cmd.Output()
	if e != nil {
		__rebarVersionErrorCached = errors.WithCause(ErrCallRebar3Command, e)
		return "", __rebarVersionErrorCached
	}
	__rebarVersionCached = string(data)
	return __rebarVersionCached, nil
}

func EvaluateRebar3Tree(ctx context.Context, dir string) ([]depNode, error) {
	var logger = logctx.Use(ctx)
	var cmd = exec.Command("rebar3", "tree")
	cmd.Dir = dir
	logger.Sugar().Infof("Execute command: %s at %s", cmd.String(), dir)
	rebarLogger := logpipe.New(logger, "rebar3")
	defer rebarLogger.Close()
	var buf = &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(buf, rebarLogger)
	cmd.Stderr = rebarLogger
	if e := cmd.Start(); e != nil {
		return nil, errors.WithCause(ErrCallRebar3Command, e)
	}
	if e := cmd.Wait(); e != nil {
		return nil, errors.WithCause(ErrCallRebar3Command, e)
	}
	return parseRebar3TreeOutput(buf.String()), nil
}
