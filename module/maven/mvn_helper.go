package maven

import (
	"github.com/ztrue/shutdown"
	"io"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type MvmCmdVersionInfo struct {
	MavenVer  string `json:"maven_ver"`
	JavaVer   string `json:"java_ver"`
	RawOutput string `json:"rawOutput"`
}

func checkMvnVersion() (*MvmCmdVersionInfo, error) {
	var mavenVersion string
	var javaVersion string
	data, e := exec.Command("mvn", "--version").Output()
	if e != nil {
		return nil, e
	}
	output := string(data)
	lines := strings.Split(output, "\n")
	versionPattern := regexp.MustCompile("Apache Maven (\\d+(?:\\.[0-9A-Za-z_-]+)+)")
	javaVersionPattern := regexp.MustCompile("Java version: (\\d+(?:\\.[0-9A-Za-z_-]+)*)")
	for _, it := range lines {
		line := strings.TrimSpace(it)
		if m := versionPattern.FindStringSubmatch(line); mavenVersion == "" && m != nil {
			mavenVersion = m[1]
			continue
		}
		if m := javaVersionPattern.FindStringSubmatch(line); javaVersion == "" && m != nil {
			javaVersion = m[1]
			continue
		}
	}

	return &MvmCmdVersionInfo{
		MavenVer:  mavenVersion,
		JavaVer:   javaVersion,
		RawOutput: output,
	}, nil
}

const _MaxMvnOutputLine = 128 * 1024

func scanMvnDependency(projectDir string) (map[Coordinate][]Dependency, error) {
	cmd := exec.Command("mvn", "dependency:tree", "-DoutputType=tgf", "--batch-mode", "-Dscope=compile")
	cmd.Dir = projectDir
	shutdownKey := shutdown.Add(func() {
		logger.Warn.Println("Maven shutdown hook execute.")
		p := cmd.Process
		if p == nil {
			logger.Warn.Println("Process is nil, skip")
			return
		}
		logger.Warn.Println("Send SIGINT to Pid:", p.Pid)
		if e := p.Signal(os.Interrupt); e != nil {
			logger.Warn.Println("Send SIGINT failed, kill.", e.Error())
			e := p.Kill()
			if e != nil {
				logger.Warn.Println("Kill process failed.", e.Error())
			}
		}
	})
	defer shutdown.Remove(shutdownKey)
	mvnStdoutR, mvnStdoutW := io.Pipe()
	defer must.Close(mvnStdoutW)
	cmd.Stdout = mvnStdoutW
	var mvnResult map[Coordinate][]Dependency
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		mvnResult = parseOutput(mvnStdoutR)
		wg.Done()
	}()
	logger.Info.Println("Execute:", cmd.String())
	if e := cmd.Start(); e != nil {
		logger.Err.Println("Mvn start failed.", e.Error())
		return nil, e
	}
	if e := cmd.Wait(); e != nil {
		logger.Err.Println("Mvn terminated with err.", e.Error())
	} else {
		logger.Info.Println("Mvn terminated with no err")
	}
	// close maven output pipe
	_ = mvnStdoutW.Close()
	wg.Wait()
	logger.Debug.Println("Mvn output parse OK.")
	return mvnResult, nil
}

// checkMvnEnv 检查maven环境
func checkMvnEnv() (bool, *MvmCmdVersionInfo) {
	if os.Getenv("NO_MVN") != "" {
		logger.Err.Println("NO_MVN environment found. Skip maven scan")
		return false, nil
	}
	ver, e := checkMvnVersion()
	if e != nil {
		logger.Err.Println("Get mvn command version failed, skip maven scan.")
		return false, nil
	}
	logger.Info.Println("Mvn command version:", ver.MavenVer)
	return true, ver
}
