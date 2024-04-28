package scanerr

import (
	"context"
	list "github.com/bahlo/generic-list-go"
)

type Param struct {
	Kind    string `json:"kind"`
	Content string `json:"content"`
}

type ctxKeyType string

var key = ctxKeyType("scanErrorCtxKey")

func WithCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, list.New[Param]())
}

func Add(ctx context.Context, e Param) {
	if l, ok := ctx.Value(key).(*list.List[Param]); ok {
		l.PushBack(e)
	}
}

func GetAll(ctx context.Context) []Param {
	l, ok := ctx.Value(key).(*list.List[Param])
	if !ok {
		return make([]Param, 0)
	}
	var rs = make([]Param, 0, l.Len())
	var f = l.Front()
	for f != nil {
		rs = append(rs, f.Value)
		f = f.Next()
	}
	return rs
}

const KindMavenNotFound = "mvn_not_found"
const KindMavenFailed = "mvn_failed"
const KindBuildDisabled = "build_disabled"
