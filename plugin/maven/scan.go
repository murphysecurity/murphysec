package maven

import (
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/util/simplejson"
	"path/filepath"
)

func executeScanCmd(pomFile string) (string, error) {
	mvnCmd := util.ExecuteCmd("mvn", "dependency:tree", "-DoutputType=dot", fmt.Sprintf("--file=%s", pomFile))
	// handling abort signal
	killSignal, canceller := util.WatchKill()
	defer canceller()
	go func() {
		if <-killSignal {
			util.KillAllChild(mvnCmd.Pid())
			mvnCmd.Abort()
			output.Warn("Scanning abort")
		}
	}()
	// print err if execute failed
	if e := mvnCmd.Execute(); e != nil {
		output.Error(fmt.Sprintf("mvn command execute failed, err: %s", e.Error()))
		if es, e := mvnCmd.GetStderr(); e == nil {
			output.Error(es)
		} else {
			output.Warn("Get mvn command error output failed.")
		}
		if es, e := mvnCmd.GetStdout(); e == nil {
			output.Info("mvn command output:")
			output.Info(es)
		}
		return "", errors.Wrap(e, "Scan project failed")
	}
	if t, e := mvnCmd.GetStdout(); e != nil {
		return "", errors.Wrap(e, "Read mvn command stdout failed")
	} else {
		return t, nil
	}
}

func doScan(dir string) (*plugin_base.PackageInfo, error) {
	mvnVer, err := mavenVersion()
	if err != nil {
		return nil, err
	}
	pom := filepath.Join(dir, "pom.xml")
	if !util.IsPathExist(pom) || util.IsDir(pom) {
		return nil, fmt.Errorf("Can't find POM file: %s\n", dir)
	}
	mvnOutput, err := executeScanCmd(pom)
	if err != nil {
		return nil, err
	}
	result, err := parseMvnCommandResult(mvnOutput)
	if err != nil {
		return nil, err
	}
	j := simplejson.NewFrom(result)
	pi := &plugin_base.PackageInfo{
		PackageManager:  "maven",
		PackageFile:     pom,
		PackageFilePath: pom,
		Language:        "java",
		Dependencies:    j,
		Name:            j.Get("name").String(),
		RuntimeInfo:     simplejson.NewFrom(mvnVer),
	}
	return pi, nil
}
