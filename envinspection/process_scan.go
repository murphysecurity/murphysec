package envinspection

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type listRunningProcessExecutableFileError struct {
	e   error
	msg string
}

func (l *listRunningProcessExecutableFileError) Error() string {
	return fmt.Sprintf("list running process failed: %s, %s", l.msg, l.e.Error())
}

func (l *listRunningProcessExecutableFileError) Unwrap() error {
	return l.e
}

func listRunningProcessExecutableFile(ctx context.Context) ([]model.Dependency, error) {
	var r []model.Dependency
	var e error
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		r, e = listRunningProcessExecutableFilePosix(ctx)
	} else if runtime.GOOS == "windows" {
		r, e = listRunningProcessExecutableFileWindows(ctx)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].Name < r[j].Name
	})
	return r, e
}

func listRunningProcessExecutableFileWindows(ctx context.Context) ([]model.Dependency, error) {
	data, e := exec.Command("wmic", "process", "get", "ExecutablePath").Output()
	if e != nil {
		return nil, &listRunningProcessExecutableFileError{e, "execute wmic failed"}
	}
	var rm = make(map[string]struct{})
	for _, s := range strings.Split(string(data), "\n") {
		s = filepath.ToSlash(strings.TrimPrefix(strings.TrimSpace(s), `\\?\`))
		if s == "" || s == "ExecutablePath" {
			continue
		}
		rm[s] = struct{}{}
	}
	var r []model.Dependency
	for s := range rm {
		r = append(r, model.Dependency{Name: s})
	}
	return r, nil
}

func listRunningProcessExecutableFilePosix(ctx context.Context) ([]model.Dependency, error) {
	var rm = make(map[string]struct{})
	entries, e := os.ReadDir("/proc")
	if e != nil {
		return nil, &listRunningProcessExecutableFileError{e, "read /proc failed"}
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		s, e := os.Readlink(filepath.Join("/proc", entry.Name(), "exe"))
		if e != nil {
			continue
		}
		rm[s] = struct{}{}
	}
	var r []model.Dependency
	for s := range rm {
		r = append(r, model.Dependency{Name: s})
	}
	return r, nil
}
