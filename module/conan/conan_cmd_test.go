package conan

import (
	"context"
	"github.com/murphysecurity/murphysec/utils/must"
	"testing"
)

func TestGetConanVersion(t *testing.T) {
	t.Log(GetConanVersion(context.TODO(), must.A(LocateConan(context.TODO()))))
}
