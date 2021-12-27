//go:build idea

package cmd

import (
	"encoding/json"
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
			TaskType:           "Cli",
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
			TaskResult: map[string]interface{}{
				"package": packageInfo.Dependencies,
				"plugin": map[string]interface{}{
					"name":    p.Info().Name,
					"version": p.Info().Version,
					"runtime": packageInfo.RuntimeInfo,
				},
			},
		})
	}()
	// idea plugin output
	if err != nil {
		fmt.Println(string(must.Byte(json.Marshal(pluginOutput{ErrCode: 1, ErrMsg: err.Error()}))))
	} else {
		fmt.Println(string(must.Byte(json.Marshal(pluginOutput{ErrCode: 0, ErrMsg: "", ScanResult: report}))))
	}
	return nil
}

type pluginOutput struct {
	*api.ScanResult
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}
