package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/plugin"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/scanner"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/version"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var scanDir string

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use:              "scan",
		Run:              scanHandler,
		TraverseChildren: true,
	}
	c.PersistentFlags().StringVarP(&scanDir, "dir", "d", ".", "project root dir")
	for i := range plugin.Plugins {
		p := plugin.Plugins[i]
		pc := &cobra.Command{
			Use:              p.Info().Name,
			Short:            p.Info().ShortDescription,
			TraverseChildren: true,
			Run: func(cmd *cobra.Command, args []string) {
				if e := scanByPlugin(p, scanDir); e != nil {
					output.Error(e.Error())
					os.Exit(-1)
				}
			},
		}
		p.SetupScanCmd(pc)
		c.AddCommand(pc)
	}
	return c
}

func scanHandler(cmd *cobra.Command, args []string) {
	panic("todo: auto scan") // todo
}

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
		util.StartSpinner("", "Waiting server response...")
		defer util.StopSpinner()
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
	output.Debug(string(must.Byte(json.Marshal(report))))
	if err != nil {
		return err
	}
	scanner.PrintDetectReport(report)
	return nil
}
