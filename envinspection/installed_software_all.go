//go:build !windows

package envinspection

import (
	"context"
	"github.com/murphysecurity/murphysec/model"
)

func listInstalledSoftwareWindows(ctx context.Context) ([]model.DependencyItem, error) {
	return nil, nil
}
