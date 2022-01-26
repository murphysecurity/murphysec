package maven

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/findfile"
	"regexp"
)

type ProjectRepo struct {
	m           map[Coordinate]*PomFile
	fileMapping map[Coordinate]string
	baseDir     string
}

func (h *ProjectRepo) FetchPomFile(ctx context.Context, coordinate Coordinate) (*PomFile, error) {
	pom := h.m[coordinate]
	if pom == nil {
		return nil, ErrArtifactNotFoundInRepo
	}
	return pom, nil
}

type ModuleInfo struct {
	*PomFile
	FilePath string
}

func (h *ProjectRepo) String() string {
	return fmt.Sprintf("ProjectRepo[%v]", h.baseDir)
}

func (h *ProjectRepo) ListModuleInfo() []ModuleInfo {
	var rs []ModuleInfo
	for _, pom := range h.m {
		rs = append(rs, ModuleInfo{
			PomFile:  pom,
			FilePath: h.fileMapping[pom.Coordinate()],
		})
	}
	return rs
}

func NewProjectRepoFromDir(dir string) (*ProjectRepo, error) {
	logger.Info.Println("Scan dir for pom file:", dir)
	iter := findfile.Find(dir, findfile.Option{
		MaxDepth:    0,
		ExcludeFile: false,
		ExcludeDir:  true,
		Predication: findfile.FileNameRegexp(regexp.MustCompile("^pom\\.xml$")),
	})
	this := &ProjectRepo{m: map[Coordinate]*PomFile{}, fileMapping: map[Coordinate]string{}, baseDir: dir}
	for iter.Next() {
		if iter.Err() != nil {
			logger.Err.Println(fmt.Sprintf("Access file failed. %+v", iter.Err()))
			continue
		}
		logger.Info.Println("Found pom file:", iter.Path())
		data, e := ioutil.ReadFile(iter.Path())
		if e != nil {
			return nil, errors.Wrap(e, "read pom file failed")
		}
		pom, e := NewPomFileFromData(data)
		if e != nil {
			return nil, e
		}
		logger.Info.Println("Read pom:", pom.Coordinate().String())
		this.m[pom.Coordinate()] = pom
		this.fileMapping[pom.Coordinate()] = iter.Path()
	}
	return this, nil
}
