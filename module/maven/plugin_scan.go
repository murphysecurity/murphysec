package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/vifraa/gopom"
	"go.uber.org/zap"
	"io/fs"
	"path/filepath"
	"time"
)

func ScanDepsByPluginCommand(ctx context.Context, projectDir string, mvnCmdInfo *MvnCommandInfo) (*DepsMap, error) {
	var logger = logctx.Use(ctx)
	var profiles, e = findPomProfiles(ctx, filepath.Join(projectDir, "pom.xml"))
	if e != nil {
		logger.Warn("Error during find pom profiles", zap.Error(e))
	} else {
		logger.Sugar().Infof("Found %d profiles", len(profiles))
	}
	c := PluginGraphCmd{
		MavenCmdInfo: mvnCmdInfo,
		Profiles:     profiles,
		Timeout:      time.Duration(env.MvnCommandTimeout) * time.Second,
		ScanDir:      projectDir,
	}
	if e := c.RunC(ctx); e != nil {
		logger.Sugar().Error("Maven graph command execution failed", zap.Error(e))
		return nil, e
	}
	logger.Sugar().Infof("Maven graph command succeeded, collecting graph file...")
	return collectPluginResultFile(ctx, projectDir)
}

func collectPluginResultFile(ctx context.Context, projectDir string) (*DepsMap, error) {
	var logger = logctx.Use(ctx)
	var graphPaths []string
	if e := filepath.Walk(projectDir, func(path string, info fs.FileInfo, err error) error {
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
	var rs = newDepsMap()
	for _, graphPath := range graphPaths {
		logger := logger.With(zap.String("graph", graphPath))
		logger.Debug("Processing graph")
		var g PluginGraphOutput
		if e := g.ReadFromFile(graphPath); e != nil {
			logger.Error("Error during read graph file", zap.Error(e))
			continue
		}
		tree, e := g.Tree()
		if e != nil {
			logger.Error("Build deps tree failed", zap.Error(e))
			continue
		}
		relPath, e := filepath.Rel(projectDir, filepath.Dir(filepath.Dir(graphPath)))
		if e != nil {
			logger.Warn("Calculate relative path failed", zap.Error(e))
		}
		rs.put(tree.Coordinate, tree.Children, filepath.Join(relPath, "pom.xml"))
	}
	return rs, nil
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
