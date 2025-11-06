//go:build !windows

package envinspection

import (
	"context"

	"github.com/murphysecurity/murphysec/model"
)

func listInstalledSoftwareWindows(ctx context.Context) ([]model.DependencyItem, error) {
	return nil, nil
}
func listPendingPatch(ctx context.Context) []string {
	return nil
}
func getWindowsVersion() model.Component {
	return model.Component{}
}
