package mvn2

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/plugin/mvn2/pom_analyzer"
	"murphysec-cli-simple/plugin/plugin_base"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/simplejson"
	"path/filepath"
)

type Plugin struct{}

var Instance plugin_base.Plugin = &Plugin{}

func (p *Plugin) SetupScanCmd(c *cobra.Command) {
}

func doScan(dir string) (map[string]*pom_analyzer.PomFile, error) {
	_, pomMap, e := pom_analyzer.PomAnalyze(dir)
	if e != nil {
		return nil, e
	}
	mvnOutput, e := executeScanCmd(filepath.Join(dir, "pom.xml"))
	if e != nil {
		return nil, e
	}
	deps := parseMvnDepOutput(mvnOutput)
	// overwrite pom analyze dependencies by mvn cmd output
	for _, it := range deps {
		if p := pomMap[it.Id()]; p != nil {
			p.Dependencies = it.Dependencies
		}
	}
	return pomMap, nil
}

type PomInfo struct {
	RelativePath string           `json:"relative_path"`
	GroupId      string           `json:"group_id"`
	ArtifactId   string           `json:"artifact_id"`
	Version      string           `json:"version"`
	Dependencies []DependencyInfo `json:"dependencies"`
}

type DependencyInfo struct {
	GroupId      string           `json:"group_id"`
	ArtifactId   string           `json:"artifact_id"`
	Version      string           `json:"version"`
	Dependencies []DependencyInfo `json:"dependencies"`
}

func (p *Plugin) DoScan(dir string) (*plugin_base.PackageInfo, error) {
	r, e := doScan(dir)
	if e != nil {
		return nil, e
	}

	rs := make([]PomInfo, 0)
	idBlackList := map[string]bool{}
	for _, it := range r {
		idBlackList[it.Id()] = true
	}
	for _, pom := range r {
		relativePath, e := filepath.Rel(dir, pom.Path)
		if e != nil {
			return nil, e
		}
		p := PomInfo{
			RelativePath: relativePath,
			GroupId:      pom.GroupId,
			ArtifactId:   pom.ArtifactId,
			Version:      pom.Version,
			Dependencies: _dependencyInfo(pom.Dependencies, idBlackList),
		}
		rs = append(rs, p)
	}
	mvnVer, _ := mavenVersion()
	packageInfo := plugin_base.PackageInfo{
		PackageManager:  "maven",
		PackageFile:     "pom.xml",
		PackageFilePath: filepath.Join(dir, "pom.xml"),
		Language:        "java",
		Dependencies:    simplejson.NewFrom(rs),
		Name:            "POM",
		RuntimeInfo:     simplejson.NewFrom(mvnVer),
	}
	return &packageInfo, nil
}

func _dependencyInfo(deps []*pom_analyzer.Dependency, blackList map[string]bool) []DependencyInfo {
	rs := make([]DependencyInfo, 0)

	for _, it := range deps {
		d := DependencyInfo{
			GroupId:      it.GroupId,
			ArtifactId:   it.ArtifactId,
			Version:      it.Version,
			Dependencies: []DependencyInfo{},
		}
		if !blackList[it.Id()] {
			d.Dependencies = _dependencyInfo(it.Dependencies, blackList)
		}
		rs = append(rs, d)
	}
	return rs
}

func (p *Plugin) MatchPath(dir string) bool {
	f := filepath.Join(must.String(filepath.Abs(dir)), "pom.xml")
	return util.IsPathExist(f) && !util.IsDir(f)
}
func (_ *Plugin) Info() *plugin_base.PluginInfo {
	return &plugin_base.PluginInfo{Name: "mvn2", ShortDescription: "for maven package"}
}
