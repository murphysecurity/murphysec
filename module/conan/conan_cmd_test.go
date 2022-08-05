package conan

import (
	"context"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"testing"
)

func TestGetConanVersion(t *testing.T) {
	_, e := exec.LookPath("conan")
	if os.Getenv("CI") != "" && errors.Is(e, exec.ErrNotFound) {
		t.Skip("Conan not found in CI environment, test skipped.")
		return
	}
	t.Log(GetConanVersion(context.TODO(), must.A(LocateConan(context.TODO()))))
}
