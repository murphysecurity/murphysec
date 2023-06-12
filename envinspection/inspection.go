package envinspection

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/conf"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/version"
	"github.com/murphysecurity/murphysec/view"
	"os"
	"runtime"
	"strings"
)

func InspectEnv(ctx context.Context, projectName string) error {
	var LOG = utils.UseLogger(ctx).Sugar()
	task := model.CreateScanTask("", model.TaskKindHostEnv, model.TaskTypeCli)
	ui := task.UI()
	ctx = model.WithScanTask(ctx, task)
	if e := createTaskApi(task, projectName); e != nil {
		return e
	}
	LOG.Infof("Task created, task id: %s", task.TaskId)

	var packageManager = "unmanaged"
	if s, ok := processByRule(version.OsName()); ok {
		packageManager = s
	}

	module2 := api.VoModule{Name: "Software Installed", PackageManager: model.PackageManagerType(packageManager)}
	var pkgs []model.Dependency
	var e error

	if runtime.GOOS == "windows" {
		pkgs, e = listInstalledSoftwareWindows(ctx)
		if e == nil && len(pkgs) > 0 {
			LOG.Warnf("Windows installed software inspection succeeded, total %d items", len(pkgs))
			module2.Dependencies = append(module2.Dependencies, pkgs...)
		} else if e != nil {
			LOG.Warnf("Windows installed software: %s", e.Error())
		}

	} else {
		pkgs, e = inspectDpkgPackage(ctx)
		if e == nil && len(pkgs) > 0 {
			LOG.Warnf("dpkg inspection succeeded, total %d items", len(pkgs))
			module2.Dependencies = append(module2.Dependencies, pkgs...)
		} else if e != nil {
			LOG.Warnf("dpkg inspection error: %s", e.Error())
		}

		pkgs, e = inspectRPMPackage(ctx)
		if e == nil && len(pkgs) > 0 {
			LOG.Warnf("RPM inspection succeeded, total %d items", len(pkgs))
			module2.Dependencies = append(module2.Dependencies, pkgs...)
		} else if e != nil {
			LOG.Warnf("RPM inspection error: %s", e.Error())
		}
	}

	if e := api.SendDetect(&api.SendDetectRequest{
		TaskInfo: task.TaskId,
		ApiToken: conf.APIToken(),
		Modules: []api.VoModule{
			module2,
			{
				Name: "OperatingSystem",
				Dependencies: []model.Dependency{
					{version.OsName(), "", nil},
				},
				PackageManager: model.PackageManagerType(packageManager),
			}},
	}); e != nil {
		LOG.Errorf("Submitting data failed: %s", e.Error())
		if e != nil {
			view.SubmitError(ui, e)
		}
		return e
	}

	if e := api.StartCheckTaskType(task.TaskId, task.Kind); e != nil {
		LOG.Errorf("Start checking failed: %s", e.Error())
		if e != nil {
			view.SubmitError(ui, e)
		}
		return e
	}

	ui.Display(display.MsgInfo, "信息提交成功")

	if r, e := api.QueryResult(task.TaskId); e != nil {
		ui.Display(display.MsgError, "获取检测结果失败: "+e.Error())
		return e
	} else {
		view.DisplayScanResultSummary(ui, r.DependenciesCount, r.IssuesCompsCount)
	}
	return nil
}

func createTaskApi(task *model.ScanTask, defaultName string) error {
	ui := task.UI()
	hostname, e := os.Hostname()
	name := defaultName
	if name == "" {
		name = fmt.Sprintf("HostEnv/%s(%s)", hostname, utils.GetOutBoundIP())
	}

	r, e := api.CreateTask(&api.CreateTaskRequest{
		CliVersion:    version.Version(),
		TaskKind:      task.Kind,
		TaskType:      task.TaskType,
		UserAgent:     version.UserAgent(),
		CmdLine:       strings.Join(os.Args, " "),
		ApiToken:      conf.APIToken(),
		ProjectName:   name,
		ProjectType:   model.ProjectTypeLocal,
		TargetAbsPath: name, // workaround
	})
	if errors.Is(e, api.ErrTlsRequest) {
		view.TLSAlert(ui, e)
		return e
	}
	if errors.Is(e, api.ErrTokenInvalid) {
		view.TokenInvalid(ui)
		return e
	}
	if e != nil {
		view.GetScanResultFailed(ui, e)
		return e
	}
	task.TaskId = r.TaskInfo
	return nil
}
