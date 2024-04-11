//go:build windows

package envinspection

import (
	"context"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows/registry"
	"testing"
)

func TestInstalledSoftwareListSubKeys(t *testing.T) {
	s, e := listSubKeys(context.TODO(), registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
	assert.NoError(t, e)
	t.Log(s)
}
