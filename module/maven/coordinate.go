package maven

import (
	"regexp"
	"strings"
)

type Coordinate struct {
	GroupId    string `json:"group_id"`
	ArtifactId string `json:"artifact_id"`
	Version    string `json:"version"`
}

var rb = regexp.MustCompile(`\s`)

func (c Coordinate) IsSnapshotVersion() bool {
	return strings.HasSuffix(c.Version, "-SNAPSHOT")
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

func (c Coordinate) Compare(o Coordinate) int {
	if n := strings.Compare(c.GroupId, o.GroupId); n != 0 {
		return n
	}
	if n := strings.Compare(c.ArtifactId, o.ArtifactId); n != 0 {
		return n
	}
	if n := strings.Compare(c.Version, o.Version); n != 0 {
		return n
	}
	return 0
}
