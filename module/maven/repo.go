package maven

import (
	"context"
	"murphysec-cli-simple/utils/semerr"
)

var ErrArtifactNotFoundInRepo = semerr.New("not found in repo")
var ErrRepoAuthRequired = semerr.New("repo auth required")
var ErrShouldRetry = semerr.New("retryable")

func DefaultMavenRepo() []Repo {
	return []Repo{MustNewHttpRepo("https://repo1.maven.org/maven2/")}
}

type Repo interface {
	FetchPomFile(ctx context.Context, coordinate Coordinate) (*PomFile, error)
}
