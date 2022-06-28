package inspector

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/model"
	"github.com/pkg/errors"
)

func createTaskC(ctx context.Context) (e error) {
	scanTask := model.UseScanTask(ctx)
	ui := scanTask.UI()
	ui.UpdateStatus(display.StatusRunning, "正在创建扫描任务，请稍候······")
	defer ui.ClearStatus()
	e = createTaskApi(ctx)
	if errors.Is(e, api.ErrTokenInvalid) {
		ui.Display(display.MsgError, "任务创建失败，Token 无效")
	} else {
		ui.Display(display.MsgError, fmt.Sprintf("任务创建失败：%s", e.Error()))
	}
	return
}
