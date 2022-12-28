package maven

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func newLocalRemote(basePath string) M2Remote {
	return &localRemote{basePath: basePath}
}

type localRemote struct {
	basePath string
}

func (l *localRemote) GetPath(_ context.Context, path string) ([]byte, error) {
	data, e := os.ReadFile(filepath.Join(l.basePath, path))
	if e != nil {
		return nil, ErrRemoteNoResource
	}
	return data, nil
}

func (l *localRemote) String() string {
	return "LocalFetcher[basePath=" + l.basePath + "]"
}

var _ fmt.Stringer = (*localRemote)(nil)
var _ M2Remote = (*localRemote)(nil)
