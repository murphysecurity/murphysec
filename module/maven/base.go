package maven

import (
	"github.com/pkg/errors"
	"github.com/vifraa/gopom"
	"murphysec-cli-simple/module/base"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/semerr"
	"path/filepath"
	"regexp"
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

func (c Coordinate) Normalize() Coordinate {
	return Coordinate{
		GroupId:    rb.ReplaceAllString(c.GroupId, ""),
		ArtifactId: rb.ReplaceAllString(c.ArtifactId, ""),
		Version:    rb.ReplaceAllString(c.Version, ""),
	}
}

func (c Coordinate) HasVersion() bool {
	return c.Normalize().Version != ""
}
func (c Coordinate) Name() string {
	//goland:noinspection GoAssignmentToReceiver
	c = c.Normalize()
	return c.GroupId + ":" + c.ArtifactId
}

func (c Coordinate) String() string {
	//goland:noinspection GoAssignmentToReceiver
	c = c.Normalize()
	if c.Version == "" {
		return c.GroupId + ":" + c.ArtifactId
	}
	return c.GroupId + ":" + c.ArtifactId + ":" + c.Version
}

var rb = regexp.MustCompile("\\s")

func (c Coordinate) IsBad() bool {
	//goland:noinspection GoAssignmentToReceiver
	c = c.Normalize()
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
	//goland:noinspection GoAssignmentToReceiver
	c = c.Normalize()
	return c.GroupId != "" && c.ArtifactId != "" && c.Version != "" && !c.IsBad()
}

var ErrInvalidCoordinate = errors.New("invalid coordinate")
var ErrArtifactNotFound = errors.New("artifact not found")
var ErrParsePomFailed = semerr.New("Parse pom failed.")

type Repo interface {
	Fetch(coordinate Coordinate) (*gopom.Project, error)
}
