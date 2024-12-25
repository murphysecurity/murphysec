package buildout

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

func base(ctx context.Context, path string, result map[string]string) error {
	var log = logctx.Use(ctx).Sugar()
	e := findVersionsFile(ctx, path, result)
	if e != nil {
		return e
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
	var extends []string
	if path == "" {
		return nil
	}
	if filepath.Dir(path) == NowPath {
		return nil
	}
	if strings.Contains(path, "http://") || strings.Contains(path, "https://") {
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
		extends, e = parseBuildoutCfgFile(ctx, path, result)
		if e != nil {
			return e
		}
	}
	for _, j := range extends {
		if j != "" {
			log.Debug("find extends", zap.String("path", j))
			e = findVersionsFile(ctx, j, result)
			if e != nil {
				return e
			}
		} else {
			e = NoCurrentDirectoryCfg(ctx, NowPath, j, result)
			if e != nil {
				return e
			}
		}
	}

	return nil
}
func findVersionsFile(ctx context.Context, path string, result map[string]string) error {
	var log = logctx.Use(ctx).Sugar()
	var extends []string
	var e error
	// 如果是远程链接 则读取
	if strings.Contains(path, "http://") || strings.Contains(path, "https://") {
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
		// 如果不是远程链接 则尝试打开读取
		extends, e = parseBuildoutCfgFile(ctx, path, result)
		if e != nil {
			return e
		}
	}
	if len(extends) != 0 {
		for _, j := range extends {
			if j != "" {
				log.Debug("find extends", zap.String("path", j))
				e = findVersionsFile(ctx, j, result)
				if e != nil {
					log.Error("find file error:", zap.Error(e))
					continue
				}
			}
		}
	}
	return nil
}
func parseBuildoutBytes(ctx context.Context, by []byte, result map[string]string) ([]string, error) {
	var log = logctx.Use(ctx).Sugar()
	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowPythonMultilineValues: true,
	}, by)
	if err != nil {
		log.Error("Fail to read file: ", zap.Error(err))
		return nil, err
	}
	for _, section := range cfg.Sections() {
		if section.Name() == "version" || section.Name() == "dependencies" || section.Name() == "versions" {
			for _, key := range section.Keys() {
				if key.Name() != "" && key.Value() != "" {
					log.Debug("buildout bytes :", zap.String(key.Name(), key.Value()))
					result[key.Name()] = key.Value()
				}
			}
		}
	}
	var resultStrings []string
	extends := cfg.Section("buildout").Key("extends").Strings("\n")
	if len(extends) == 0 {
		return nil, nil
	}
	for _, j := range extends {
		if j != "" {
			resultStrings = append(resultStrings, j)
		}
	}
	return resultStrings, nil
}
func parseBuildoutCfgFile(ctx context.Context, path string, result map[string]string) ([]string, error) {
	var log = logctx.Use(ctx).Sugar()
	by, err := os.ReadFile(path)
	if err != nil {
		log.Error("read file failed", zap.Error(err))
		return nil, err
	}
	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowPythonMultilineValues: true,
	}, by)
	if err != nil {
		log.Error("Fail to read file: ", zap.Error(err))
		return nil, err
	}

	for _, section := range cfg.Sections() {
		if section.Name() == "version" || section.Name() == "dependencies" || section.Name() == "versions" {
			for _, key := range section.Keys() {
				if key.Name() != "" && key.Value() != "" {
					log.Debug("from path:", zap.String(path, key.Name()))
					result[key.Name()] = key.Value()
				}
			}
		}
	}
	var resultStrings []string
	extends := cfg.Section("buildout").Key("extends").Strings("\n")
	if len(extends) == 0 {
		return nil, nil
	}
	for _, j := range extends {
		if j != "" && !filepath.IsAbs(j) && !strings.Contains(j, "http://") && !strings.Contains(j, "https://") {
			resultStrings = append(resultStrings, filepath.Join(filepath.Dir(path), j))
		} else {
			resultStrings = append(resultStrings, j)
		}
	}
	return resultStrings, nil
}
