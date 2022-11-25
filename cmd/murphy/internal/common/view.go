package common

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/ui"
)

func displayTokenNotSet(ctx context.Context) {
	ui.Use(ctx).Display(ui.MsgError, "您似乎没有设置令牌，请使用 auth login 命令，或者 --token 参数设置一个访问令牌")
}

func displayGetTokenErr(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "读取访问令牌失败："+e.Error())
}

func displayInitializeFailed(ctx context.Context, e error) {
	ui.Use(ctx).Display(ui.MsgError, "初始化失败："+e.Error())
}
