package maven

import (
	"bufio"
	"fmt"
	"github.com/ztrue/shutdown"
	"io"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var modulePattern = regexp.MustCompile("digraph +\\\"(.+?):(.+?):.+?:(.+?)\\\" .*\\{")
var depPattern = regexp.MustCompile("\\\"([^:\\\"]+):([^:\\\"]+):(?:[^:\\\"]+):([^:\\\"]+)(?::([^:\\\"]+))?\\\"\\s*->\\s*\\\"([^:\\\"]+):([^:\\\"]+):(?:[^:\\\"]+):([^:\\\"]+)(?::([^:\\\"]+))?\\\"")

const _MaxLineSize = 128 * 1024

func parseOutput(reader io.Reader) map[Coordinate][]Dependency {
	collection := map[Coordinate]map[Coordinate][]Coordinate{}

	input := bufio.NewScanner(reader)
	input.Split(bufio.ScanLines)
	input.Buffer(make([]byte, _MaxLineSize), _MaxLineSize)

	var currentModule Coordinate

	for input.Scan() {
		line := input.Text()
		if m := modulePattern.FindStringSubmatch(line); m != nil {
			// module matched
			currentModule = Coordinate{
				GroupId:    m[1],
				ArtifactId: m[2],
				Version:    m[3],
			}
			continue
		}
		if m := depPattern.FindStringSubmatch(line); m != nil {
			if collection[currentModule] == nil {
				collection[currentModule] = map[Coordinate][]Coordinate{}
			}
			if !utils.InStringSlice([]string{"compile", "runtime", ""}, m[4]) ||
				!utils.InStringSlice([]string{"compile", "runtime", ""}, m[8]) {
				logger.Debug.Println("Skip line", line)
				continue
			}
			left := Coordinate{
				GroupId:    m[1],
				ArtifactId: m[2],
				Version:    m[3],
			}
			collection[currentModule][left] = append(collection[currentModule][left], Coordinate{
				GroupId:    m[5],
				ArtifactId: m[6],
				Version:    m[7],
			})
			continue
		}
		logger.Debug.Println("Unrecognized line:", line)
	}
	logger.Debug.Println("Total modules:", len(collection))
	var count int
	for _, it := range collection {
		for _, it := range it {
			count += len(it)
		}
	}
	logger.Debug.Println("Total items:", count)
	graphs := map[Coordinate][]Dependency{}
	for module, it := range collection {
		graphs[module] = _conv(module, it, nil).Children
	}
	return graphs
}
func _conv(root Coordinate, m map[Coordinate][]Coordinate, visited []Coordinate) Dependency {
	for i := range visited {
		if visited[i] == root {
			var names []string
			for _, it := range visited {
				names = append(names, it.String())
			}
			logger.Warn.Println("Circular dependency:", strings.Join(names, " -> "))
			return Dependency{Coordinate: root}
		}
	}
	rootDep := Dependency{Coordinate: root}
	if len(visited) < 6 {
		for _, it := range m[root] {
			rootDep.Children = append(rootDep.Children, _conv(it, m, append(visited, root)))
		}
	}
	return rootDep
}

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
	cmd := exec.Command("mvn", "dependency:tree", "-DoutputType=dot", fmt.Sprintf("--file=%s/pom.xml", projectDir))
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
		return nil, e
	} else {
		logger.Info.Println("Mvn terminated with no err")
	}
	// close maven output pipe
	_ = mvnStdoutW.Close()
	wg.Wait()
	logger.Debug.Println("Mvn output parse OK.")
	return mvnResult, nil
}
