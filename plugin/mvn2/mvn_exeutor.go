package mvn2

import (
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/plugin/mvn2/pom_analyzer"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/output"
	"regexp"
	"strings"
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

func mavenVersion() (*RuntimeMavenVersion, error) {
	c := util.ExecuteCmd("mvn", "--version")
	killSig, canceller := util.WatchKill()
	defer canceller()
	go func() {
		if <-killSig {
			util.KillAllChild(c.Pid())
			c.Abort()
		}
	}()
	if e := c.Execute(); e != nil {
		if s, e := c.GetStderr(); e != nil {
			output.Warn(fmt.Sprintf("Get error out failed: %s", e.Error()))
		} else {
			output.Warn(s)
		}
		return nil, errors.Wrap(e, "Get maven version failed")
	}
	if t, e := c.GetStdout(); e == nil {
		return parseMvnVerCommandResult(t), nil
	} else {
		return nil, errors.Wrap(e, "Read maven stdout failed")
	}
}

type RuntimeMavenVersion struct {
	MvnVersion  string `json:"mvn_version"`
	JavaVersion string `json:"java_version"`
	RuntimeOs   string `json:"runtime_os"`
}

func parseMvnVerCommandResult(cmdResult string) *RuntimeMavenVersion {
	lines := strings.Split(cmdResult, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	rs := RuntimeMavenVersion{}
	for _, it := range lines {
		switch {
		case strings.HasPrefix(it, "Apache Maven"):
			rs.MvnVersion = it
		case strings.HasPrefix(it, "Java version"):
			rs.JavaVersion = it
		case strings.HasPrefix(it, "Os name"):
			rs.RuntimeOs = it
		}
	}
	return &rs
}

func parseMvnDepOutput(input string) []*pom_analyzer.Dependency {
	p := regexp.MustCompile("digraph +\\\"(.+?):(.+?):.+?:(.+?)\\\" .*\\{")
	p2 := regexp.MustCompile("\\\"(?:(.+?):(.+?):.+?:(.+?))[:\\\"].*?->\\s+\\\"(?:(.+?):(.+?):.+?:(.+?))[\\\":]")
	lines := strings.Split(input, "\n")

	var modules []*pom_analyzer.Dependency
	var currentModule *pom_analyzer.Dependency
	var depsMap map[string]*pom_analyzer.Dependency
	for _, line := range lines {
		if m := p.FindStringSubmatch(line); m != nil {
			// if module headline
			currentModule = &pom_analyzer.Dependency{
				GroupId:    m[1],
				ArtifactId: m[2],
				Version:    m[3],
			}
			modules = append(modules, currentModule)
			depsMap = map[string]*pom_analyzer.Dependency{currentModule.Id(): currentModule}
		} else if m := p2.FindStringSubmatch(line); m != nil && currentModule != nil {
			// if currentModule != nil && is dependency line
			leftId := fmt.Sprintf("%s:%s:%s", m[1], m[2], m[3])
			rightId := fmt.Sprintf("%s:%s:%s", m[4], m[5], m[6])
			// create dependency object if not exists
			if depsMap[leftId] == nil {
				depsMap[leftId] = &pom_analyzer.Dependency{
					GroupId:    m[1],
					ArtifactId: m[2],
					Version:    m[3],
				}
			}
			if depsMap[rightId] == nil {
				depsMap[rightId] = &pom_analyzer.Dependency{
					GroupId:    m[4],
					ArtifactId: m[5],
					Version:    m[6],
				}
			}
			// associate them
			depsMap[leftId].Dependencies = append(depsMap[leftId].Dependencies, depsMap[rightId])
		}
	}
	return modules
}
