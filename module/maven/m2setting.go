package maven

import (
	"bytes"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
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
		LocalRepoPath: filepath.Join(must.String(homedir.Dir()), ".m2", "repository"),
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
		opt.LocalRepoPath = strings.ReplaceAll(n.InnerText(), "${user.home}", must.String(homedir.Dir()))
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
	if hd, e := homedir.Dir(); e != nil {
		return nil, e
	} else {
		p := filepath.Join(hd, ".m2", "settings.xml")
		return os.ReadFile(p)
	}
}

func readMvnInstallPathSettingsFile() ([]byte, error) {
	mvnBin := locateMvnInstallPath()
	if mvnBin == "" {
		return nil, errors.New("mvn binary not found")
	}
	p := filepath.Join(filepath.Base(filepath.Base(mvnBin)), "conf", "settings.xml")
	return os.ReadFile(p)
}
