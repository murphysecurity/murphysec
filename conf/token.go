package conf

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const tokenPath = "~/.murphysec/token"

// APITokenCliOverride 用于覆盖API TOKEN，为空串时无效
var APITokenCliOverride string

// _APITokenEnvOverride 用于环境变量覆盖 API TOKEN，为空串时无效
var _APITokenEnvOverride = func() string { return os.Getenv("API_TOKEN") }()

// tokenReader read token file from config directory, the file will only be read once time.
var tokenReader = func() func() string {
	o := sync.Once{}
	var t = ""
	return func() string {
		o.Do(func() {
			dir, e := homedir.Expand(tokenPath)
			if e != nil {
				logger.Debug.Printf("Cannot get home path, ignore")
				return
			}
			logger.Debug.Printf("Read token from: %s", dir)
			data, e := ioutil.ReadFile(dir)
			if e != nil {
				logger.Debug.Printf("Read fail, ignore")
				return
			}
			t = strings.TrimSpace(string(data))
		})
		return t
	}
}()

func ReadTokenFile() (t string, e error) {
	defer func() {
		if e != nil {
			e = errors.Wrap(e, "Read token failed")
		}
	}()
	dir, e := homedir.Expand(tokenPath)
	if e != nil {
		logger.Err.Println("Expand home dir failed,", e.Error())
		return "", e
	}
	logger.Debug.Println("Read token file at:", dir)
	data, e := ioutil.ReadFile(dir)
	if e != nil {
		logger.Debug.Println("Read fail failed", e.Error())
		return "", e
	}
	return strings.TrimSpace(string(data)), nil
}

// APIToken returns API token
func APIToken() string {
	if len(strings.TrimSpace(APITokenCliOverride)) != 0 {
		logger.Debug.Printf("Use API token from cli argument")
		return APITokenCliOverride
	}
	if len(strings.TrimSpace(_APITokenEnvOverride)) != 0 {
		logger.Debug.Printf("Use API token from env")
		return _APITokenEnvOverride
	}
	logger.Debug.Printf("Use API token from config file")
	return tokenReader()
}

// StoreToken store specified token to user local config
func StoreToken(token string) error {
	path, err := homedir.Expand(tokenPath)
	if err != nil {
		return errors.Wrap(err, "Can't get your home dir.")
	}
	if e := os.MkdirAll(filepath.Dir(path), 0777); e != nil {
		return errors.Wrap(e, "Create config dir failed.")
	}
	if e := ioutil.WriteFile(path, []byte(token), 0600); e != nil {
		return errors.Wrap(e, "Write token file failed.")
	}
	return nil
}

// TokenFileNotFound will be returned when the token file not found
var TokenFileNotFound = errors.New("TokenFileNotFound")

// RemoveToken will delete local token file, return TokenFileNotFound if there is no token.
func RemoveToken() error {
	path, err := homedir.Expand(tokenPath)
	if err != nil {
		return errors.Wrap(err, "Can't get your home dir.")
	}
	if stat, e := os.Stat(path); e != nil || stat.IsDir() {
		return TokenFileNotFound
	}
	if e := os.Remove(path); e != nil {
		return errors.Wrap(e, "Delete token file failed.")
	}
	return nil
}
