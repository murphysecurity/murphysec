//go:build windows

package envinspection

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
)

type listSubKeysError struct {
	e       error
	message string
}

func (l *listSubKeysError) Unwrap() error {
	return l.e
}

func (l *listSubKeysError) Error() string {
	return fmt.Sprintf("listSubkeys: %s", l.e)
}

func listSubKeys(ctx context.Context, key registry.Key, path string) ([]string, error) {
	k, e := registry.OpenKey(key, path, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if e != nil {
		return nil, &listSubKeysError{e, "open key failed"}
	}
	defer k.Close()
	ki, e := k.Stat()
	if e != nil {
		return nil, &listSubKeysError{e, "get stat failed"}
	}
	if ki.SubKeyCount < 0 {
		return nil, &listSubKeysError{e, "SubKeyCount < 0"}
	}
	skList, e := k.ReadSubKeyNames(int(ki.SubKeyCount))
	if e != nil {
		return nil, &listSubKeysError{e, "read subkey names failed"}
	}
	var r = make([]string, 0, len(skList))
	for _, s := range skList {
		r = append(r, filepath.Join(path, s))
	}
	return r, nil
}

func listInstalledSoftwareWindows(ctx context.Context) ([]model.Dependency, error) {
	paths, e := listSubKeys(ctx, registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
	if e != nil {
		return nil, e
	}
	var r []model.Dependency
	for _, p := range paths {
		k, e := registry.OpenKey(registry.LOCAL_MACHINE, p, registry.READ)
		if e != nil {
			continue
		}
		displayName, _, _ := k.GetStringValue("DisplayName")
		displayVersion, _, _ := k.GetStringValue("DisplayVersion")
		_ = k.Close()
		if displayName == "" {
			continue
		}
		r = append(r, model.Dependency{
			Name:    displayName,
			Version: displayVersion,
		})
	}
	return r, nil
}
