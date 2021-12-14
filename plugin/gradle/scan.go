package gradle

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/util/simplejson"
	"os"
	"path/filepath"
	"strings"
)

func scanDir(abortSignal chan struct{}, dir string) (string, error) {
	doneSignal := make(chan struct{})
	defer close(doneSignal)
	// collect information
	gradleCmd := getGradleCmd(dir)
	gradleFile := detectGradleFile(dir)
	if gradleFile == "" {
		return "", errors.New("No gradle build file found, supported: build.gradle, build.gradle.kts")
	}
	version, err := detectGradleVersion(gradleCmd)
	if err != nil {
		return "", errors.Wrap(err, "Detect gradle version failed")
	}
	// prepare scan script
	scanScriptPath, cleanTemp, err := tempScanScript()
	if err != nil {
		output.Error(err.Error())
		return "", err
	}
	defer cleanTemp()
	// execute scan script
	output.Debug(fmt.Sprintf("Use gradle path: %s, version: %s", gradleCmd, version.Version))
	cmd := util.ExecuteCmd(gradleCmd, "getDepsJson", "-q", "--build-file="+gradleFile, "--no-daemon", "-Dorg.gradle.parallel=", "-Dorg.gradle.console=plain", "-I", scanScriptPath)
	go func() {
		select {
		case <-abortSignal:
			util.KillAllChild(cmd.Pid())
			cmd.Abort()
			output.Warn("Scan abort.")
		case <-doneSignal:
		}
	}()
	if err := cmd.Execute(); err != nil {
		output.Error(fmt.Sprintf("Execute scan script failed, %v", err))
		es, _ := cmd.GetStderr()
		output.Error(es)
		return "", err
	} else {
		es, e := cmd.GetStdout()
		if e != nil {
			output.Error(fmt.Sprintf("Read gradle output failed, %s", e.Error()))
			return "", e
		}
		return es, nil
	}
}

func parseGradleScanCmdResult(cmdResult string) (interface{}, error) {
	depsInfo := strings.Trim(cmdResult, "GetDepsJson:")
	var j = simplejson.New()
	if e := json.Unmarshal([]byte(depsInfo), &j); e != nil {
		output.Error("parse scan result failed")
		output.Error(e.Error())
		return nil, e
	} else {
		output.Debug("scan result parsed")
		output.Debug(j.MarshalString())
		return j, nil
	}
}

func tempScanScript() (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "murphysec-")
	if err != nil {
		return "", nil, errors.Wrap(err, "Create temp dir failed")
	}
	output.Debug(fmt.Sprintf("Make temp dir succeed, %s", tempDir))
	p := filepath.Join(tempDir, "murphysec-scan.gradle")
	err = ioutil.WriteFile(p, []byte(initScriptContent), 644)
	if err != nil {
		return "", nil, errors.Wrap(err, "Write temp file failed")
	}
	output.Debug("Write temp file succeed")
	cleanup := func() {
		output.Debug(fmt.Sprintf("cleanup temp scan script: %s", tempDir))
		e := os.RemoveAll(tempDir)
		if e != nil {
			output.Warn(fmt.Sprintf("Failed, %v", e))
		}
		output.Debug("Succeed")
	}
	return p, cleanup, nil
}
