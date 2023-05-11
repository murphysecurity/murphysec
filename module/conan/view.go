package conan

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/ui"
)

func badConanView(ctx context.Context) {
	u := ui.Use(ctx)
	u.Display(ui.MsgWarn, "识别到您的环境中 conan 无法正常运行，可能会导致检测结果不完整或失败，访问 https://murphysec.com/docs/faqs/quick-start-for-beginners/programming-language-supported.html 了解详情")
}

func printConanError(ctx context.Context, e *conanError) {
	u := ui.Use(ctx)
	for _, s := range e.ErrorMultiLine() {
		u.Display(ui.MsgWarn, s)
	}
}
