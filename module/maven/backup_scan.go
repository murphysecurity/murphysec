package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
	"net/url"
	"path/filepath"
)

func BackupResolve(ctx context.Context, projectDir string) (*DepsMap, error) {
	var rs = newDepsMap()
	var logger = utils.UseLogger(ctx)
	logger.Sugar().Infof("Backup scan: %s", projectDir)
	resolver, e := prepareResolver(ctx)
	if e != nil {
		return nil, e
	}
	poms, e := ReadLocalProject(ctx, projectDir)
	if e != nil {
		return nil, e
	}
	logger.Sugar().Infof("Found %d pom file", len(poms))
	for _, pom := range poms {
		resolver.pomCache.add(pom)
	}
	for _, pom := range poms {
		coordinate := pom.Coordinate()
		if !coordinate.Complete() {
			logger.Warn("Incomplete coordinate, skip", zap.Any("coordinate", coordinate))
			continue
		}
		logger.Sugar().Infof("Build dependency tree: %s", pom.Coordinate())
		tree := BuildDepTree(ctx, resolver, coordinate)
		relPath, e := filepath.Rel(projectDir, pom.Path)
		if e != nil {
			logger.Warn("Calculate relative-path failed", zap.Error(e), zap.Any("path", pom.Path))
		}
		rs.put(coordinate, tree.Children, relPath)
	}
	return rs, nil
}

func prepareResolver(ctx context.Context) (*PomResolver, error) {
	logger := utils.UseLogger(ctx)
	userConfig, e := GetMvnConfig(ctx)
	if e != nil {
		return nil, e
	}
	logger.Sugar().Infof("User maven config: %s", userConfig.String())
	resolver := NewPomResolver(ctx)
	if userConfig.Repo != "" {
		resolver.AddRepo(NewLocalRepo(userConfig.Repo))
		logger.Sugar().Debugf("Add local repo: %s", userConfig.Repo)
	}
	for _, remote := range userConfig.Remotes {
		u, e := url.Parse(remote)
		if e != nil {
			logger.Warn("Parse url failed", zap.Error(e), zap.String("remote", remote))
			continue
		}
		httpRepo := NewHttpRepo(ctx, *u)
		resolver.AddRepo(httpRepo)
		logger.Sugar().Debugf("Add http repo: %s", httpRepo)
	}
	return resolver, nil
}
