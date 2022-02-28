package maven

import (
	"github.com/pkg/errors"
	"github.com/vifraa/gopom"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/semerr"
	"path/filepath"
	"strings"
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

func (c Coordinate) IsBad() bool {
	if strings.HasPrefix(c.GroupId, "${") ||
		strings.HasPrefix(c.ArtifactId, "${") ||
		strings.HasPrefix(c.Version, "${") ||
		strings.HasPrefix(c.Version, "[") ||
		strings.HasPrefix(c.Version, "(") {
		return true
	}
	return false
}

func (c Coordinate) Complete() bool {
	return c.GroupId != "" && c.ArtifactId != "" && c.Version != "" && !c.IsBad()
}

var ErrInvalidCoordinate = errors.New("invalid coordinate")
var ErrArtifactNotFound = errors.New("artifact not found")
var ErrParsePomFailed = semerr.New("Parse pom failed.")

type Repo interface {
	Fetch(coordinate Coordinate) (*gopom.Project, error)
}
