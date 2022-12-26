package maven

import (
	"context"
	list "github.com/bahlo/generic-list-go"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

func ReadLocalProject(ctx context.Context, dir string) ([]*UnresolvedPom, error) {
	logger := logctx.Use(ctx)

	// Workaround for Maven CI Friendly Versions: https://maven.apache.org/maven-ci-friendly.html
	var revisionMap = map[string]string{}

	var moduleQ = list.New[string]()
	moduleQ.PushBack(dir)

	var projectPomList []*UnresolvedPom

	var visitedPath = make(map[string]bool)
	for moduleQ.Len() > 0 {
		current := moduleQ.Front().Value
		moduleQ.Remove(moduleQ.Front())

		if visitedPath[current] {
			continue
		}
		visitedPath[current] = true

		pomPath := filepath.Join(current, "pom.xml")
		pom, e := readPomFile(ctx, pomPath)
		if e != nil {
			logger.Warn("Read pom failed", zap.String("path", current), zap.Error(e))
			continue
		}

		if pom.Properties.Entries != nil {
			if v := pom.Properties.Entries["revision"]; v != "" && !strings.Contains(v, "${") {
				revisionMap[pom.GroupID+pom.ArtifactID] = v
				if pom.Version == "${revision}" {
					pom.Version = v
				}
			}
		}
		if pom.Version != "" && !strings.Contains(pom.Version, "${") {
			revisionMap[pom.GroupID+pom.ArtifactID] = pom.Version
		}
		if pom.Parent.Version == "${revision}" {
			if v := revisionMap[pom.Parent.GroupID+pom.Parent.ArtifactID]; v != "" {
				pom.Parent.Version = v
			}
		}

		for _, module := range pom.Modules {
			moduleQ.PushBack(filepath.Join(current, module))
		}

		for _, profile := range pom.Profiles {
			for _, module := range profile.Modules {
				moduleQ.PushBack(filepath.Join(current, module))
			}
		}

		projectPomList = append(projectPomList, &UnresolvedPom{pom, pomPath})
	}

	return projectPomList, nil
}
