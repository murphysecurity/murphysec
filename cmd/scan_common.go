//go:build !embedding && !idea

package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/scanner"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/util/spin_util"
	"murphysec-cli-simple/version"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func scanByPlugin(p plugin_base.Plugin, dir string) error {
	startTime := time.Now()
	dir = must.String(filepath.Abs(dir))
	output.Info(fmt.Sprintf("Scan dir: %s", dir))
	if !p.MatchPath(dir) {
		return errors.New(fmt.Sprintf("The project can't be processed by plugin %s.", p.Info().Name))
	}
	packageInfo, err := p.DoScan(dir)
	if err != nil {
		return err
	}
	output.Info("Finish package scan, collecting project information...")
	projectInfo := scanner.GetProjectInfo(dir)
	output.Info("Project information collected.")
	must.True(projectInfo != nil)
	scanEndTime := time.Now()
	// do report
	report, err := func() (*api.ScanResult, error) {
		spin_util.StartSpinner("", "Waiting server response...")
		defer spin_util.StopSpinner()
		return api.Report(&api.ScanRequestBody{
			CliVersion:         version.Version(),
			TaskStatus:         1,
			TaskFailureReason:  "",
			TaskType:           "CLI",
			OsType:             runtime.GOOS,
			CmdLine:            strings.Join(os.Args, " "),
			Plugin:             p.Info().Name,
			TaskConsumeTime:    int(scanEndTime.Sub(startTime).Seconds()),
			ApiToken:           conf.APIToken(),
			TaskStartTimestamp: int(startTime.Unix()),
			ProjectType:        projectInfo.ProjectType,
			ProjectName:        projectInfo.ProjectName,
			GitRemoteUrl:       projectInfo.GitRemoteUrl,
			GitBranch:          projectInfo.GitBranch,
			TargetPath:         dir,
			TargetAbsPath:      dir,
			PackageManager:     packageInfo.PackageManager,
			PackageFile:        packageInfo.PackageFile,
			PackageFilePath:    packageInfo.PackageFilePath,
			Language:           packageInfo.Language,
			Dependencies:       packageInfo.Dependencies,
			RuntimeInfo:        packageInfo.RuntimeInfo,
			TaskResult:         struct{}{},
		})
	}()
	if err != nil {
		return err
	}
	scanner.PrintDetectReport(report)
	return nil
}
