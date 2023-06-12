//go:build !windows

package envinspection

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
)

func listInstalledSoftwareWindows(ctx context.Context) ([]model.Dependency, error) {
	return nil, nil
}
