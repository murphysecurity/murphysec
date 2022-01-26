package maven

import (
	"context"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
	"strings"
)

type LocalRepo struct {
	basePath string
}

var _LocalRepoInstance = &LocalRepo{
	basePath: filepath.Join(must.String(homedir.Dir()), ".m2", "repository"),
}

func NewLocalRepo(path string) *LocalRepo {
	if path != "" {
		_LocalRepoInstance.basePath = path
	}
	return _LocalRepoInstance
}

func (h *LocalRepo) FetchPomFile(ctx context.Context, coordinate Coordinate) (*PomFile, error) {
	pathSeg := []string{mavenHome(), "repository"}
	pathSeg = append(pathSeg, strings.Split(coordinate.GroupId, ".")...)
	pathSeg = append(pathSeg, coordinate.ArtifactId, coordinate.Version, fmt.Sprintf("%s-%s.pom", coordinate.ArtifactId, coordinate.Version))
	pomData, e := ioutil.ReadFile(filepath.Join(pathSeg...))
	if e != nil {
		return nil, ErrArtifactNotFoundInRepo
	}
	p, e := NewPomFileFromData(pomData)
	if e != nil {
		return nil, errors.Wrap(e, "invalid pom file")
	}
	if p.Coordinate() != coordinate {
		return nil, ErrArtifactNotFoundInRepo
	}
	logger.Debug.Println("Local repo load:", p.Coordinate())
	return p, nil
}

func (h *LocalRepo) String() string {
	return fmt.Sprintf("LocalRepo[%v]", mavenHome())
}

func mavenHome() string {
	var m2home string
	if s := os.Getenv("M2_HOME"); s != "" {
		m2home = s
	}
	if m2home == "" {
		m2home = filepath.Join(must.String(homedir.Dir()), ".m2")
	}
	return m2home
}
