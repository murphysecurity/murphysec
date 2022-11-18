package envinspection

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/version"
	"github.com/murphysecurity/murphysec/view"
	"os"
	"strings"
)

func InspectEnv(ctx context.Context) error {
	var LOG = utils.UseLogger(ctx).Sugar()
	task := model.CreateScanTask("", model.TaskKindHostEnv, model.TaskTypeCli)
	ui := task.UI()
	ctx = model.WithScanTask(ctx, task)
	if e := createTaskApi(task); e != nil {
		return e
	}
	LOG.Infof("Task created, task id: %s", task.TaskId)

	var kernelVer = readLinuxKernelVersion()
	osRelease := readOsRelease()
	osId := osRelease["ID"]
	osVersion := osRelease["VERSION"]

	module2 := api.VoModule{Name: "Software Installed", PackageManager: "Unmanaged"}
	pkgs, e := inspectDpkgPackage(ctx)

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

	if e := api.SendDetect(&api.SendDetectRequest{
		TaskInfo: task.TaskId,
		ApiToken: "", // todo: token empty
		Modules: []api.VoModule{
			module2,
			{
				Name: "OperatingSystem",
				Dependencies: []model.Dependency{
					{"kernel", kernelVer, nil},
					{osId, osVersion, nil},
				},
				PackageManager: "Unmanaged",
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

func createTaskApi(task *model.ScanTask) error {
	ui := task.UI()
	hostname, e := os.Hostname()
	name := fmt.Sprintf("HostEnv/%s(%s)", hostname, utils.GetOutBoundIP())
	if e != nil {
		return errors.WithCause(ErrGetHostname, e)
	}

	r, e := api.CreateTask(&api.CreateTaskRequest{
		CliVersion:    version.Version(),
		TaskKind:      task.Kind,
		TaskType:      task.TaskType,
		UserAgent:     version.UserAgent(),
		CmdLine:       strings.Join(os.Args, " "),
		ApiToken:      "", // todo: empty token
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
