package maven

import (
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vifraa/gopom"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type LocalRepo struct {
	baseDir string
}

func NewLocalRepo(s string) *LocalRepo {
	return &LocalRepo{baseDir: s}
}

func (l *LocalRepo) String() string {
	return fmt.Sprintf("LocalRepo[%s]", l.baseDir)
}

func (l *LocalRepo) Fetch(coordinate Coordinate) (*gopom.Project, error) {
	if !coordinate.Complete() {
		return nil, ErrInvalidCoordinate
	}
	p := filepath.Join(l.baseDir)
	for _, s := range strings.Split(coordinate.GroupId, ".") {
		p = filepath.Join(p, s)
	}
	p = filepath.Join(p, coordinate.ArtifactId, coordinate.Version, coordinate.ArtifactId+"-"+coordinate.Version+".pom")
	if info, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return nil, ErrArtifactNotFound
	} else {
		if info.IsDir() {
			return nil, errors.New("it's a directory")
		}
	}
	data, e := ioutil.ReadFile(p)
	if e != nil {
		return nil, errors.Wrap(e, "Read local pom file failed")
	}
	var proj gopom.Project
	if e := xml.Unmarshal(data, &proj); e != nil {
		return nil, ErrParsePomFailed.Decorate(e)
	}
	return &proj, nil
}
