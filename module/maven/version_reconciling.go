package maven

import (
	"context"
	"github.com/murphysecurity/murphysec/utils"
	"go.uber.org/zap"
)

func VersionReconciling(ctx context.Context, root *Dependency) Dependency {
	c := newReconcilingCtx(ctx)
	c._visit(*root)
	return c._assign(*root)
}

type reconcilingCtx struct {
	logger     *zap.Logger
	versionMap map[string]string
}

func newReconcilingCtx(ctx context.Context) *reconcilingCtx {
	return &reconcilingCtx{
		logger:     utils.UseLogger(ctx),
		versionMap: map[string]string{},
	}
}

func (r *reconcilingCtx) _visit(root Dependency) {
	k := root.GroupId + ":" + root.ArtifactId
	if r.versionMap[k] == "" && root.Version != "" {
		r.versionMap[k] = root.Version
	}
	for idx := range root.Children {
		r._visit(root.Children[idx])
	}
}

func (r *reconcilingCtx) _assign(root Dependency) Dependency {
	k := root.GroupId + ":" + root.ArtifactId
	root.Version = r.versionMap[k]
	for i := range root.Children {
		root.Children[i] = r._assign(root.Children[i])
	}
	return root
}
