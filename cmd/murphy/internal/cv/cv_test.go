package cv

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/ui"
	"testing"
	"time"
)

func TestView(t *testing.T) {
	ctx := ui.With(context.TODO(), ui.CLI{})
	DisplayScanning(ctx)
	DisplaySubtaskCreated(ctx, "项目名称", "任务名称", "123", "aa", "456")
	time.Sleep(time.Second / 4)
	DisplayTLSNotice(ctx)
}
