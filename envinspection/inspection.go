package envinspection

import (
	"context"
	"errors"
	"fmt"
	"github.com/iseki0/osname"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/scanerr"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
)

func InspectEnv(ctx context.Context) error {
	var LOG = logctx.Use(ctx).Sugar()
	task := model.UseScanTask(ctx)
	if task == nil {
		panic("task == nil")
	}

	var packageManager = "unmanaged"
	var osn, _ = osname.OsName()
	if runtime.GOOS == "linux" {
		packageManager = getOsInfo()
	}
	if s, ok := processByRule(osn); ok {
		packageManager = s
	}
	var module = model.Module{
		ModuleName:     "InstalledSoftware",
		PackageManager: packageManager,
		Dependencies:   nil,
		ModulePath:     "/InstalledSoftware", // never be empty, workaround for the platform issue
	}

	var scanFunc []func(ctx context.Context) ([]model.DependencyItem, error)
	if runtime.GOOS == "windows1" {
		scanFunc = append(scanFunc, listInstalledSoftwareWindows /*listRunningProcessExecutableFileWindows*/)
	} else {
		scanFunc = append(scanFunc, listDpkgPackage, listRPMPackage /*listRunningProcessExecutableFilePosix*/)
	}
	var foundCmd = false
	for _, f := range scanFunc {
		pkgs, e := f(ctx)
		var fn = strings.TrimPrefix(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), "github.com/murphysecurity/murphysec/envinspection.")
		if e == nil && len(pkgs) > 0 {
			module.Dependencies = append(module.Dependencies, pkgs...)
			LOG.Infof("inspection succeeded(%s), total %d items", fn, len(pkgs))
			continue
		} else {
			LOG.Warnf("Software inspection error(%s): %s, ", fn, e)
		}
		if errors.Is(e, exec.ErrNotFound) {
			continue
		}
		foundCmd = true
		var cError cError
		if errors.As(e, &cError) {
			if cError.Content == "" {
				cError.Content = "(no stderr output)"
			}
			scanerr.Add(ctx, scanerr.Param{
				Kind:    "env_inspection_error",
				Content: cError.Content,
			})
		}
	}
	if !foundCmd {
		scanerr.Add(ctx, scanerr.Param{
			Kind:    "env_inspection_unsupported",
			Content: fmt.Sprintf("no command found for %s", runtime.GOOS),
		})
	}
	for i := range module.Dependencies {
		module.Dependencies[i].IsOnline.SetOnline(false)
		module.Dependencies[i].IsDirectDependency = true
		module.Dependencies[i].EcoRepo.Repository = packageManager
	}
	task.Modules = append(task.Modules, module)

	return nil
}
