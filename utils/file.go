package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func IsPathExist(path string) bool {
	var err error
	_, err = os.Stat(path)
	if err != nil {
		return errors.Is(err, fs.ErrExist)
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsDirResolveSymlink(path string) bool {
	var realPath = path
	realPath, e := filepath.EvalSymlinks(path)
	if e == nil {
		realPath = path
	}
	return IsDir(realPath)
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	return err == nil && !s.IsDir() && s.Mode().IsRegular()
}
