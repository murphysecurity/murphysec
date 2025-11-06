//go:build windows

package envinspection

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/repeale/fp-go"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
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

func getWindowsVersion() model.Component {
	return model.Component{
		CompName:    "Windows",
		CompVersion: fmt.Sprintf("%d.%d.%d", windows.RtlGetVersion().MajorVersion, windows.RtlGetVersion().MinorVersion, windows.RtlGetVersion().BuildNumber),
	}
}

func listInstalledSoftwareWindows(ctx context.Context) ([]model.DependencyItem, error) {
	var rKeys = []registry.Key{registry.CURRENT_USER, registry.LOCAL_MACHINE}
	var r []model.DependencyItem
	for _, rKey := range rKeys {
		paths, e := listSubKeys(ctx, rKey, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
		if e != nil {
			return nil, e
		}
		for _, p := range paths {
			k, e := registry.OpenKey(rKey, p, registry.READ)
			if e != nil {
				continue
			}
			displayName, _, _ := k.GetStringValue("DisplayName")
			displayVersion, _, _ := k.GetStringValue("DisplayVersion")
			_ = k.Close()
			if displayName == "" {
				continue
			}
			r = append(r, model.DependencyItem{
				Component: model.Component{
					CompName:    displayName,
					CompVersion: displayVersion,
				},
			})
		}
	}
	return r, nil
}

func listPendingPatch(ctx context.Context) []string {
	var logger = logctx.Use(ctx).Sugar()
	data, e := exec.CommandContext(ctx, "wmic", "qfe", "get", "HotFixID").Output()
	if e != nil {
		logger.Warnf("wmic qfe get HotFixID failed: %s", e)
	}
	var isKB = func(s string) bool { return strings.HasPrefix(s, "KB") }
	return fp.Pipe2(fp.Map(strings.TrimSpace), fp.Filter(isKB))(strings.Split(string(data), "\n"))
}
