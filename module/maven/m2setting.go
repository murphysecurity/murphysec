package maven

import (
	"bytes"
	"github.com/antchfx/xmlquery"
	"github.com/mitchellh/go-homedir"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"os"
	"path/filepath"
)

func ReadM2SettingMirror() *M2Setting {
	// todo: refactor
	rs := &M2Setting{
		Mirrors:  nil,
		RepoPath: filepath.Join(must.String(homedir.Dir()), ".m2", "repository"),
	}
	settingPath := filepath.Join(must.String(homedir.Dir()), ".m2", "settings.xml")
	if t := os.Getenv("MPS_CLI_M2_SETTINGS"); t != "" {
		settingPath = t
	}
	logger.Info.Println("Read maven settings.xml from", settingPath)
	data, e := os.ReadFile(settingPath)
	if e != nil {
		logger.Err.Println("Read settings failed.", e.Error())
		return nil
	}
	node, e := xmlquery.Parse(bytes.NewReader(data))
	if e != nil {
		logger.Err.Println("Parse m2 settings failed.", e.Error())
	}
	for _, it := range xmlquery.Find(node, "/settings/mirrors/mirror") {
		url := xmlquery.FindOne(it, "url")
		if url == nil {
			continue
		}
		rs.Mirrors = append(rs.Mirrors, url.InnerText())
	}
	if n := xmlquery.FindOne(node, "/settings/localRepository"); n != nil {
		rs.RepoPath = n.InnerText()
	}
	logger.Info.Println("m2 info", rs)
	return rs
}

type M2Setting struct {
	Mirrors  []string
	RepoPath string
}
