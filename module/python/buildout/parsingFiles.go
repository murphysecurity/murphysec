package buildout

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func base(ctx context.Context, path string, result map[string]string) error {
	var log = logctx.Use(ctx).Sugar()
	e := findVersionsFile(ctx, path, result)
	if e == nil {
		return nil
	}
	pattern := filepath.Join(filepath.Dir(path), "*.cfg")
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Error("glob failed", zap.Error(e))
		return err
	}
	for _, j := range files {
		if j == "buildout.cfg" {
			continue
		}
		if err := NoCurrentDirectoryCfg(ctx, filepath.Dir(j), j, result); err != nil {
			return err
		}
	}
	return nil
}
func NoCurrentDirectoryCfg(ctx context.Context, NowPath string, path string, result map[string]string) error {
	var log = logctx.Use(ctx).Sugar()
	var e error
	var extends = ""
	if path == "" {
		return nil
	}
	if filepath.Dir(path) == NowPath {
		return nil
	}
	if strings.Contains(path, "http") {
		resp, err := http.Get(path)
		if err != nil {
			log.Error("http get failed", zap.Error(err))
			return err
		}
		defer resp.Body.Close()
		by, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("read body failed", zap.Error(err))
			return err
		}
		extends, e = parseBuildoutBytes(ctx, by, result)
		if e != nil {
			return e
		}
	} else {
		// 如果不是远程连接 则尝试打开读取
		extends, e = parseBuildoutCfgFile(ctx, path, result)
		if e != nil {
			return e
		}
	}
	if extends != "" {
		log.Debug("find extends", zap.String("path", extends))
		return findVersionsFile(ctx, extends, result)
	}
	return NoCurrentDirectoryCfg(ctx, NowPath, extends, result)
}
func findVersionsFile(ctx context.Context, path string, result map[string]string) error {
	var log = logctx.Use(ctx).Sugar()
	var extends string
	var e error
	// 如果事远程链接 则读取
	if strings.Contains(path, "http") {
		resp, err := http.Get(path)
		if err != nil {
			log.Error("http get failed", zap.Error(err))
			return err
		}
		defer resp.Body.Close()
		by, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("read body failed", zap.Error(err))
			return err
		}
		extends, e = parseBuildoutBytes(ctx, by, result)
		if e != nil {
			return e
		}
	} else {
		// 如果不是远程连接 则尝试打开读取
		extends, e = parseBuildoutCfgFile(ctx, path, result)
		if e != nil {
			return e
		}
	}
	if extends != "" {
		log.Debug("find extends", zap.String("path", extends))
		return findVersionsFile(ctx, extends, result)
	}

	return nil
}
func parseBuildoutBytes(ctx context.Context, by []byte, result map[string]string) (string, error) {
	var log = logctx.Use(ctx).Sugar()
	cfg, err := ini.Load(by)
	if err != nil {
		log.Error("Fail to read file: ", zap.Error(err))
		return "", err
	}
	for _, section := range cfg.Sections() {
		if section.Name() == "version" || section.Name() == "dependencies" || section.Name() == "versions" {
			for _, key := range section.Keys() {
				if key.Name() != "" && key.Value() != "" {
					result[key.Name()] = key.Value()
				}
			}
		}
	}
	extends := cfg.Section("buildout").Key("extends").String()
	return extends, nil
}
func parseBuildoutCfgFile(ctx context.Context, path string, result map[string]string) (string, error) {
	var log = logctx.Use(ctx).Sugar()
	by, err := os.ReadFile(path)
	if err != nil {
		log.Error("read file failed", zap.Error(err))
		return "", err
	}
	cfg, err := ini.Load(by)
	if err != nil {
		log.Error("Fail to read file: ", zap.Error(err))
		return "", err
	}
	for _, section := range cfg.Sections() {
		if section.Name() == "version" || section.Name() == "dependencies" || section.Name() == "versions" {
			for _, key := range section.Keys() {
				if key.Name() != "" && key.Value() != "" {
					result[key.Name()] = key.Value()
				}
			}
		}
	}
	extends := cfg.Section("buildout").Key("extends").String()
	if extends == "" {
		return "", nil
	}
	if filepath.IsAbs(extends) {
		return extends, nil
	} else {
		return filepath.Join(filepath.Dir(path), extends), nil
	}
}
