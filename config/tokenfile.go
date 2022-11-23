package config

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	DefaultTokenFile = "~/.murphysec/token"
	TokenFileMaxSize = 1 * 1024
)

var tokenPattern = regexp.MustCompile(`^[\w._+\/-]+$`)

func ReadTokenFile(ctx context.Context) (token string, e error) {
	logger := logctx.Use(ctx).Sugar()
	token, e = _readTokenFile(ctx)
	if e != nil {
		logger.Errorf("Read token: %v", e)
		e = ErrNoToken
	}
	return
}
func _readTokenFile(ctx context.Context) (token string, e error) {
	var (
		data          []byte
		s             string
		f             *os.File
		tokenFilePath string
	)
	tokenFilePath, e = homedir.Expand(DefaultTokenFile)
	if e != nil {
		return
	}
	f, e = os.Open(tokenFilePath)
	if e != nil {
		return
	}
	defer func() { _ = f.Close() }()

	data, e = io.ReadAll(io.LimitReader(f, TokenFileMaxSize))
	if e != nil {
		return
	}

	s = strings.TrimSpace(string(data))
	if !tokenPattern.MatchString(s) {
		e = _ErrTokenFileReadFailed
		return
	}
	token = s
	return
}

func WriteLocalTokenFile(ctx context.Context, token string) error {
	var (
		e      error
		fp     string
		logger = logctx.Use(ctx).Sugar()
	)
	if !tokenPattern.MatchString(token) {
		return ErrBadToken
	}
	logger.Infof("update local token")
	fp, e = homedir.Expand(DefaultTokenFile)
	if e != nil {
		logger.Errorf("get token path failed: %v", e)
		return e
	}
	e = os.MkdirAll(filepath.Dir(fp), 0755)
	if e != nil {
		logger.Errorf("MkdirAll: %v", e)
	}
	e = os.WriteFile(fp, []byte(token), 0755)
	if e != nil {
		logger.Error("write token file failed")
	}
	logger.Infof("complete")
	return nil
}

func RemoveTokenFile(ctx context.Context) error {
	var (
		fp, e = homedir.Expand(DefaultTokenFile)
	)
	if e != nil {
		return e
	}
	if !utils.IsPathExist(fp) {
		return nil
	}
	if !utils.IsFile(fp) {
		return _ErrTokenFileNotAFile
	}
	return os.Remove(fp)
}
