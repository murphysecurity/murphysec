package maven

import (
	"context"
	"encoding/xml"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/vifraa/gopom"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ScanMvnDeps(ctx context.Context, mvnCmdInfo *MvnCommandInfo) (map[Coordinate][]Dependency, error) {
	logger := utils.UseLogger(ctx)
	scanTask := model.UseInspectorTask(ctx)
	startDir := scanTask.ScanDir
	var profiles, e = findPomProfiles(ctx, filepath.Join(startDir, "pom.xml"))
	if e != nil {
		logger.Sugar().Warnf("Error during find pom profiles: %v", e)
	}
	logger.Sugar().Infof("Found profiles: %s", strings.Join(profiles, ","))
	c := MvnGraphCmdArgs{
		Path:     mvnCmdInfo.Path,
		Profiles: profiles,
		Timeout:  time.Duration(env.MvnCommandTimeout) * time.Second,
		ScanDir:  startDir,
	}
	if e := c.Execute(ctx); e != nil {
		logger.Sugar().Error("Maven graph command execution failed", zap.Error(e))
		return nil, e
	}
	logger.Sugar().Infof("Maven graph command succeeded, collecting graph file...")
	var graphPaths []string
	if e := filepath.Walk(startDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}
		if info.Name() == "dependency-graph.json" {
			logger.Sugar().Debugf("Found graph file: %s", path)
			graphPaths = append(graphPaths, path)
		}
		return nil
	}); e != nil {
		logger.Warn("Error during graph collection.", zap.Error(e))
	}
	depsMap := map[Coordinate][]Dependency{}
	for _, graphPath := range graphPaths {
		logger = logger.With(zap.String("graph", graphPath))
		logger.Debug("Processing graph")
		coordinate := readCoordinate(filepath.Dir(filepath.Dir(graphPath)))
		if coordinate == nil {
			logger.Error("Read coordinate failed")
			continue
		}
		var g *dependencyGraph
		if e := g.ReadFromFile(graphPath); e != nil {
			logger.Error("Error during read graph file", zap.Error(e))
			continue
		}
		tree, e := g.Tree()
		if e != nil {
			logger.Error("Build deps tree failed", zap.Error(e))
			continue
		}
		depsMap[*coordinate] = tree
	}
	return depsMap, nil
}

func findPomProfiles(ctx context.Context, pomPath string) (profiles []string, err error) {
	project, e := gopom.Parse(pomPath)
	if e != nil {
		return nil, e
	}
	for _, profile := range project.Profiles {
		profiles = append(profiles, profile.ID)
	}
	return
}

func readCoordinate(dir string) *Coordinate {
	data, e := os.ReadFile(filepath.Join(dir, "pom.xml"))
	if e != nil {
		return nil
	}
	var p gopom.Project
	if e := xml.Unmarshal(data, &p); e != nil {
		return nil
	}
	c := &Coordinate{
		GroupId:    p.GroupID,
		ArtifactId: p.ArtifactID,
		Version:    p.Version,
	}
	if c.GroupId == "" {
		c.GroupId = p.Parent.GroupID
	}
	if c.ArtifactId == "" {
		c.ArtifactId = p.Parent.ArtifactID
	}
	if c.Version == "" {
		c.Version = p.Parent.Version
	}
	return c
}
