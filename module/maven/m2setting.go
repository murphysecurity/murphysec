package maven

import (
	"bytes"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/mitchellh/go-homedir"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

type MvnOption struct {
	LocalRepoPath string
	Remote        []string
}

func (o MvnOption) String() string {
	return fmt.Sprintf("LocalRepo: %s, remotes: %s", o.LocalRepoPath, strings.Join(o.Remote, ","))
}

func DefaultMvnOption() MvnOption {
	return MvnOption{
		LocalRepoPath: filepath.Join(must.A(homedir.Dir()), ".m2", "repository"),
		Remote:        []string{"https://repo1.maven.org/maven2/"},
	}
}

func ReadMvnOption() MvnOption {
	var data []byte
	var e error
	data, e = readUserHomeM2Settings()
	if e != nil {
		logger.Info.Println("Read user home maven settings.xml failed.", e.Error())
		data, e = readMvnInstallPathSettingsFile()
		if e != nil {
			logger.Info.Println("Read maven install settings.xml failed.", e.Error())
		}
	}
	node, e := xmlquery.Parse(bytes.NewReader(data))
	if e != nil {
		logger.Info.Println("Parse m2 settings failed.", e.Error())
		return DefaultMvnOption()
	}
	opt := DefaultMvnOption()
	opt.Remote = nil
	for _, it := range xmlquery.Find(node, "/settings/mirrors/mirror") {
		url := xmlquery.FindOne(it, "url")
		if url == nil {
			continue
		}
		opt.Remote = append(opt.Remote, url.InnerText())
	}
	opt.Remote = append(opt.Remote, "https://repo1.maven.org/maven2/")
	if n := xmlquery.FindOne(node, "/settings/localRepository"); n != nil {
		opt.LocalRepoPath = strings.ReplaceAll(n.InnerText(), "${user.home}", must.A(homedir.Dir()))
	}
	logger.Info.Println("maven option", opt)
	return opt
}

type M2Setting struct {
	Mirrors  []string
	RepoPath string
}

func locateMvnInstallPath() string {
	for _, it := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		fp, e := filepath.EvalSymlinks(filepath.Join(it, "mvn"))
		if e != nil {
			continue
		}
		return fp
	}
	return ""
}

func readUserHomeM2Settings() ([]byte, error) {
	var baseDir = os.Getenv("M2_HOME")
	if baseDir == "" {
		var e error
		baseDir, e = homedir.Dir()
		if e != nil {
			return nil, e
		}
		baseDir = filepath.Join(baseDir, ".m2")
	}
	return os.ReadFile(filepath.Join(baseDir, "settings.xml"))
}

func readMvnInstallPathSettingsFile() ([]byte, error) {
	mvnBin := locateMvnInstallPath()
	if mvnBin == "" {
		return nil, errors.New("mvn binary not found")
	}
	p := filepath.Join(filepath.Dir(filepath.Dir(mvnBin)), "conf", "settings.xml")
	return os.ReadFile(p)
}
