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
	"sync"
)

const (
	GlobalConfDir    = "~/.murphysec"
	TokenFilePath    = GlobalConfDir + "/token"
	TokenFileMaxSize = 1 * 1024
)

var CliTokenOverride string

var tokenPattern = regexp.MustCompile(`^[\w._+\/-]+$`)

func GetToken(ctx context.Context) (string, error) {
	initToken(ctx)
	return _token, _tokenErr
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
	fp, e = homedir.Expand(TokenFilePath)
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

var _token string
var _tokenErr error
var _loadTokenOnce sync.Once

func initToken(ctx context.Context) {
	_loadTokenOnce.Do(func() {
		var (
			logger    = logctx.Use(ctx).Sugar().Named("tokenInit")
			envr      = os.Getenv("API_TOKEN")
			fileToken string
		)
		logger.Debugf("begin")
		defer func() { logger.Infof("initToken complete, %v", _tokenErr) }()
		if CliTokenOverride != "" {
			_token = CliTokenOverride
			logger.Infof("argument overriding")
			return
		}
		if envr != "" {
			_token = envr
			logger.Infof("environment variable overriding")
			return
		}

		fileToken, _tokenErr = readTokenFile(ctx)
		if _tokenErr != nil {
			logger.Infof("read token file failed: %v", _tokenErr)
			_tokenErr = nil
		} else {
			logger.Infof("use token file")
			_token = fileToken
			return
		}
		_tokenErr = ErrNoToken
		return
	})
}

func ReadTokenFile(ctx context.Context) (string, error) {
	return readTokenFile(ctx)
}

func RemoveTokenFile(ctx context.Context) error {
	var (
		fp, e = homedir.Expand(TokenFilePath)
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

func readTokenFile(ctx context.Context) (token string, e error) {
	var (
		data          []byte
		s             string
		f             *os.File
		tokenFilePath string
	)
	tokenFilePath, e = homedir.Expand(TokenFilePath)
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
