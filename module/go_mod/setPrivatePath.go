package go_mod

import (
	"os"

	"go.uber.org/zap"
)

func setPrivatePath(privatePath string, logger *zap.Logger) error {
	err := os.Setenv("GOPRIVATE", privatePath)
	if err != nil {
		logger.Error("Error setting GOPRIVATE environment variable", zap.Error(err))
		return err
	}
	return nil
}
func setProxyPath(proxyPath string, logger *zap.Logger) error {
	err := os.Setenv("GOPROXY", proxyPath)
	if err != nil {
		logger.Error("Error setting GOPROXY environment variable", zap.Error(err))
		return err
	}
	return nil
}
