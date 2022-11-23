package ui

import (
	"context"
)

type uiKeyType struct{}

var uiKey uiKeyType

func With(ctx context.Context, ui UI) context.Context {
	if ui == nil {
		panic("ui is nil")
	}
	return context.WithValue(ctx, uiKey, ui)
}

func Use(ctx context.Context) UI {
	d, _ := ctx.Value(uiKey).(UI)
	return d
}
