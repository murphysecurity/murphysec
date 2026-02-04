package envinspection

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/iseki0/osname"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/scanerr"
)

func InspectEnv(ctx context.Context, scanProcess bool, windowsPatchScanTimeout time.Duration) error {
	task := model.UseScanTask(ctx)
	if task == nil {
		panic("task == nil")
	}

	var packageManager = "unmanaged"
	var osn, _ = osname.OsName()
	if runtime.GOOS == "linux" {
		packageManager = getOsInfo()
	} else if runtime.GOOS == "windows" {
		version := getWindowsVersion()
		var m = model.Module{
			ModuleName:   "5ec239b6-715c-4d36-a3b8-a5a629b898a9",
			Dependencies: []model.DependencyItem{{Component: version}},
			Patches:      listPendingPatch(ctx, windowsPatchScanTimeout),
			ModulePath:   "/Windows",
		}
		task.Modules = append(task.Modules, m)
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

	// 获取软件包列表
	inspectInstalledSoftware(ctx, &module)

	if scanProcess {
		// 获取进程文件列表
		inspectProcessFiles(ctx)
	}

	return nil
}

// 检查已安装的软件
func inspectInstalledSoftware(ctx context.Context, module *model.Module) {
	LOG := logctx.Use(ctx).Sugar()
	task := model.UseScanTask(ctx)

	var scanFunc []func(ctx context.Context) ([]model.DependencyItem, error)
	if runtime.GOOS == "windows" {
		scanFunc = append(scanFunc, listInstalledSoftwareWindows /*listRunningProcessExecutableFileWindows*/)
	} else {
		scanFunc = append(scanFunc, listDpkgPackage, listRPMPackage /*listRunningProcessExecutableFilePosix*/)
	}
	var foundCmd = false
	for _, f := range scanFunc {
		var fn = strings.TrimPrefix(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), "github.com/murphysecurity/murphysec/envinspection.")
		pkgs, e := f(ctx)
		if e != nil {
			LOG.Warnf("Software inspection error(%s): %s, ", fn, e)
			if errors.Is(e, fs.ErrNotExist) || errors.Is(e, exec.ErrNotFound) {
				LOG.Debug("Skipping command due to missing executable: ", fn)
				continue
			}
			var pError *exec.ExitError
			if errors.As(e, &pError) {
				var stderrText = strings.TrimSpace(string(pError.Stderr))
				if stderrText == "" {
					stderrText = "(no stderr output)"
				}
				scanerr.Add(ctx, scanerr.Param{
					Kind:    "env_inspection_error",
					Content: stderrText,
				})
				continue
			}
			foundCmd = true
			continue
		}
		module.Dependencies = append(module.Dependencies, pkgs...)
		foundCmd = true
	}
	if !foundCmd {
		scanerr.Add(ctx, scanerr.Param{
			Kind:    "env_inspection_unsupported",
			Content: fmt.Sprintf("no command found for %s", runtime.GOOS),
		})
	}
	var ecosystem = runtime.GOOS
	if ecosystem != "linux" && ecosystem != "windows" {
		ecosystem = "generic"
	}
	for i := range module.Dependencies {
		module.Dependencies[i].IsOnline.SetOnline(false)
		module.Dependencies[i].DependencyRelation = model.DependencyRelationDirect
		module.Dependencies[i].EcoRepo.Repository = module.PackageManager
		module.Dependencies[i].EcoRepo.Ecosystem = ecosystem
	}
	task.Modules = append(task.Modules, *module)
}

// 检查进程文件
func inspectProcessFiles(ctx context.Context) {
	LOG := logctx.Use(ctx).Sugar()
	task := model.UseScanTask(ctx)

	// 仅在Linux系统上执行进程文件检查
	if runtime.GOOS != "linux" {
		LOG.Info("Process file inspection is only supported on Linux, skipping...")
		return
	}

	// 获取进程相关的模块列表
	processModules, err := listProcessFiles(ctx)
	if err != nil {
		LOG.Warnf("Process file inspection error: %s", err)
		scanerr.Add(ctx, scanerr.Param{
			Kind:    "process_file_inspection_error",
			Content: err.Error(),
		})
		return
	}

	// 添加进程模块到任务中
	if len(processModules) > 0 {
		task.Modules = append(task.Modules, processModules...)
		LOG.Infof("Process file inspection succeeded, found %d processes", len(processModules))
	} else {
		LOG.Info("No process files found")
	}
}
