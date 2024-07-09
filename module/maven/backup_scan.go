package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
	"net/url"
	"path/filepath"
)

func BackupResolve(ctx context.Context, projectDir string) (*DepsMap, error) {
	var rs = newDepsMap()
	var logger = logctx.Use(ctx)
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
		resolver.addPom(pom)
	}
	for _, pom := range poms {
		coordinate := pom.Coordinate()
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
	logger := logctx.Use(ctx)
	userConfig, e := GetMvnConfig(ctx)
	if e != nil {
		return nil, e
	}
	logger.Sugar().Infof("User maven config: %s", userConfig.String())
	var remotes []M2Remote

	if userConfig.Repo != "" {
		remotes = append(remotes, newLocalRemote(userConfig.Repo))
		logger.Sugar().Debugf("Add local repo: %s", userConfig.Repo)
	}
	for _, remote := range userConfig.Remotes {
		u, e := url.Parse(remote)
		if e != nil {
			logger.Warn("Parse url failed", zap.Error(e), zap.String("remote", remote))
			continue
		}
		httpRemote := newHttpRemote(*u)
		remotes = append(remotes, httpRemote)
		logger.Sugar().Debugf("Add http repo: %s", httpRemote)
	}

	return NewPomResolver(ctx, remotes), nil
}
