package maven

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type UserConfig struct {
	Remotes []string
	Repo    string
}

func (u UserConfig) String() string {
	return fmt.Sprintf("[Repo=%s, Mirrors=%s]", u.Repo, strings.Join(u.Remotes, ","))
}

func defaultUserConfig() *UserConfig {
	p, _ := homedir.Expand(".m2/repository")
	u := &UserConfig{
		Remotes: nil,
		Repo:    p,
	}
	if env.MavenCentral != "" {
		u.Remotes = append(u.Remotes, env.MavenCentral)
	}
	return u
}

func GetMvnConfig(ctx context.Context) (*UserConfig, error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	logger := logctx.Use(ctx)
	var node *xmlquery.Node
	for _, p := range mavenSettingsPaths(ctx) {
		logger.Debug("Reading maven settings", zap.String("path", p))
		if !utils.IsFile(p) {
			logger.Debug("not a file, skip")
			continue
		}
		data, e := os.ReadFile(p)
		if e != nil {
			logger.Warn("Read failed", zap.Error(e))
			continue
		}
		node, e = xmlquery.Parse(bytes.NewReader(data))
		if e != nil {
			logger.Warn("Parse failed", zap.Error(e))
			continue
		}
		break
	}
	if node == nil {
		logger.Info("No maven settings found, use default config")
		return defaultUserConfig(), nil
	}

	var uc = &UserConfig{}
	for _, it := range xmlquery.Find(node, "/settings/mirrors/mirror") {
		url := xmlquery.FindOne(it, "url")
		if url == nil {
			continue
		}
		uc.Remotes = append(uc.Remotes, url.InnerText())
	}
	if n := xmlquery.FindOne(node, "/settings/localRepository"); n != nil {
		uc.Repo = strings.ReplaceAll(n.InnerText(), "${user.home}", must.A(homedir.Dir()))
	}
	if env.MavenCentral != "" {
		uc.Remotes = append(uc.Remotes, env.MavenCentral)
	}
	return uc, nil
}

func locateMvnInstallPath(ctx context.Context) string {
	logger := logctx.Use(ctx)
	cmdPath := locateMvnCmdPath()
	if cmdPath == "" {
		return ""
	}
	fp, e := filepath.EvalSymlinks(cmdPath)
	if e != nil {
		return ""
	}
	p := filepath.Dir(filepath.Dir(fp))
	logger.Sugar().Debugf("Maven install at: %s", p)
	return p
}

func mavenSettingsPaths(ctx context.Context) (paths []string) {
	// IDEA specified path
	if env.IdeaMavenConf != "" {
		paths = append(paths, env.IdeaMavenConf)
	}
	// user path
	var homeDir = os.Getenv("M2_HOME")
	if homeDir == "" {
		if b, e := homedir.Dir(); e == nil {
			homeDir = filepath.Join(b, ".m2")
		}
	}
	if homeDir != "" {
		paths = append(paths, filepath.Join(homeDir, "settings.xml"))
	}

	// install path
	base := locateMvnInstallPath(ctx)
	var candidate = []string{"conf/settings.xml"}
	if runtime.GOOS == "darwin" {
		candidate = append(candidate, "libexec/conf/settings.xml")
	}
	for _, s := range candidate {
		paths = append(paths, filepath.Join(base, s))
	}
	return
}
