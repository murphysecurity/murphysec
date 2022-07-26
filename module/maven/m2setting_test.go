package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/utils/must"
	"testing"
)

func TestGetMvnConfig(t *testing.T) {
	t.Log(must.A(GetMvnConfig(context.TODO())))
}
