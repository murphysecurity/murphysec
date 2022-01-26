package maven

import (
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"path/filepath"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "MavenInspector@" + i.Version()
}

func (i *Inspector) Version() string {
	return "v0.0.1"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "pom.xml"))
}

func (i *Inspector) Inspect(dir string) ([]base.Module, error) {
	return ScanMavenProject(dir)
}

func (i *Inspector) PackageManagerType() base.PackageManagerType {
	return base.PMMaven
}

type Coordinate struct {
	GroupId    string `json:"group_id"`
	ArtifactId string `json:"artifact_id"`
	Version    string `json:"version"`
}

func (c Coordinate) HasVersion() bool {
	return c.Version != ""
}
func (c Coordinate) Name() string {
	return c.GroupId + ":" + c.ArtifactId
}

func (c Coordinate) String() string {
	if c.Version == "" {
		return c.GroupId + ":" + c.ArtifactId
	}
	return c.GroupId + ":" + c.ArtifactId + ":" + c.Version
}
