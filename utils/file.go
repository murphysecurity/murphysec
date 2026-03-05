package utils

import (
	"os"
	"path/filepath"
)

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsDirResolveSymlink(path string) bool {
	realPath := path
	if p, e := filepath.EvalSymlinks(path); e == nil {
		realPath = p
	}
	return IsDir(realPath)
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	return err == nil && !s.IsDir() && s.Mode().IsRegular()
}
