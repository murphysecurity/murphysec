package go_mod

import (
	"context"
	"path/filepath"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Inspector struct{}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return model.InspectorFeatureAllowNested&feature > 0
}

func (Inspector) String() string {
	return "GoMod"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, "go.mod"))
}

func (Inspector) InspectProject(ctx context.Context) error {
	logger := logctx.Use(ctx)
	if privatePath, ok := ctx.Value("privateSourceAddr").(string); ok {
		logger.Debug("Use private path", zap.String("path", privatePath))
		if err := setPrivatePath(privatePath, logger); err != nil {
			return err
		}
	}
	if proxyPath, ok := ctx.Value("proxyAddr").(string); ok {
		logger.Debug("Use proxy path", zap.String("path", proxyPath))
		if err := setProxyPath(proxyPath, logger); err != nil {
			return err
		}
	}
	if err := buildScan(ctx); err != nil {
		if err := baseScan(ctx); err != nil {
			return err
		}
	}
	return nil
}

var EcoRepo = model.EcoRepo{
	Ecosystem:  "go",
	Repository: "",
}

var _ErrGoNotFound = errors.New("go not found")
