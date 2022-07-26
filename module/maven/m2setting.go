package maven

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/env"
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
	logger := utils.UseLogger(ctx)
	var node *xmlquery.Node
	for _, p := range mavenSettingsPaths() {
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

func locateMvnInstallPath() string {
	info, e := CheckMvnCommand()
	if e != nil {
		return ""
	}
	fp, e := filepath.EvalSymlinks(info.Path)
	if e == nil {
		return filepath.Dir(filepath.Dir(fp))
	}
	return ""
}

func mavenSettingsPaths() (paths []string) {
	// user path
	var baseDir = os.Getenv("M2_HOME")
	if baseDir == "" {
		if b, e := homedir.Dir(); e == nil {
			baseDir = filepath.Join(b, ".m2")
		}
	}
	paths = append(paths, baseDir)

	// install path
	base := locateMvnInstallPath()
	var candidate = []string{"conf/settings.xml"}
	if runtime.GOOS == "darwin" {
		candidate = append(candidate, "libexec/conf/settings.xml")
	}
	for _, s := range candidate {
		paths = append(paths, filepath.Join(base, s))
	}
	return
}
