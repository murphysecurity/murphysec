//go:build !windows

package envinspection

import (
	"context"
	"time"

	"github.com/murphysecurity/murphysec/model"
)

func listInstalledSoftwareWindows(ctx context.Context) ([]model.DependencyItem, error) {
	return nil, nil
}
func listPendingPatch(ctx context.Context, windowsPatchScanTimeout time.Duration) []string {
	return nil
}
func getWindowsVersion() model.Component {
	return model.Component{}
}
