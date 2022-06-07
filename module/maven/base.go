package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/base"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/semerr"
	"github.com/pkg/errors"
	"github.com/vifraa/gopom"
	"path/filepath"
	"regexp"
	"strings"
)

type Inspector struct{}

func New() base.Inspector {
	return &Inspector{}
}

func (i *Inspector) String() string {
	return "MavenInspector"
}

func (i *Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "pom.xml"))
}

func (i *Inspector) InspectProject(ctx context.Context) error {
	task := model.UseInspectorTask(ctx)
	modules, e := ScanMavenProject(task)
	if e != nil {
		return e
	}
	for _, it := range modules {
		task.AddModule(it)
	}
	return nil
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
