package cv

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/ui"
	"github.com/murphysecurity/murphysec/model"
	"testing"
)

func TestDisplayScanResultSummary(t *testing.T) {
	var ctx = ui.With(context.TODO(), ui.CLI)
	DisplayScanResultSummary(ctx, 10, 5, 3, []model.ScanWarning{{Kind: "test_foo"}, {Kind: "bar"}})
}
