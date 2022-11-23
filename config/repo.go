package config

import (
	"context"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

const DefaultRepoConfigName = ".murphy.yml"
const RepoConfigVersion = "1.0.0"

type configFileT struct {
	Version string     `json:"version" yaml:"version"`
	IDEA    RepoConfig `json:"idea,omitempty" yaml:"idea,omitempty"`
	CLI     RepoConfig `json:"cli,omitempty" yaml:"cli,omitempty"`
}

func (c *configFileT) GetAccessType(accessType model.AccessType) RepoConfig {
	switch accessType {
	case model.AccessTypeCli:
		return c.CLI
	case model.AccessTypeIdea:
		return c.IDEA
	default:
		panic("bad accessType")
	}
}

func (c *configFileT) SetByAccessType(accessType model.AccessType, config RepoConfig) {
	switch accessType {
	case model.AccessTypeCli:
		c.CLI = config
	case model.AccessTypeIdea:
		c.IDEA = config
	default:
		panic("bad accessType")
	}
}

type RepoConfig struct {
	TaskId string `json:"task_id,omitempty" yaml:"task_id,omitempty"`
}

func ReadRepoConfig(ctx context.Context, repoPath string, accessType model.AccessType) (*RepoConfig, error) {
	if !accessType.Valid() {
		return nil, ErrRepoConfigBadAccessType
	}

	var logger = logctx.Use(ctx).Sugar()
	var cf *configFileT
	var e error
	var f *os.File
	var optionFilePath = repoConfigPath(repoPath)
	logger.Debugf("Use config: %s", optionFilePath)

	f, e = os.Open(optionFilePath)
	if e != nil {
		return nil, ErrRepoConfigNotFound
	}
	defer func() { utils.LogCloseErr(logger, "config-file", f) }()

	cf, e = readRepoConfigFromPipe(ctx, f)
	if e != nil {
		return nil, e
	}

	var selectedConfig = cf.GetAccessType(accessType)
	return &selectedConfig, nil
}

func WriteRepoConfig(ctx context.Context, repoPath string, accessType model.AccessType, cf RepoConfig) error {
	must.NotNil(ctx)

	if !accessType.Valid() {
		return ErrRepoConfigBadAccessType
	}
	var logger = logctx.Use(ctx).Sugar()
	var f *os.File
	var e error
	var optionFilePath = repoConfigPath(repoPath)
	logger.Debugf("Use config: %s", optionFilePath)

	// open file
	f, e = os.OpenFile(optionFilePath, os.O_CREATE+os.O_RDWR, 0644)
	if e != nil {
		logger.Errorf("open config file failed: %v", e)
		return e
	}
	defer func() { utils.LogCloseErr(logger, "config-file", f) }()

	// read file
	var m *configFileT
	m, e = readRepoConfigFromPipe(ctx, f)
	if e != nil {
		return e
	}

	_, e = f.Seek(0, 0)
	if e != nil {
		logger.Errorf("seek: %v", e)
		return e
	}

	m.Version = RepoConfigVersion
	m.SetByAccessType(accessType, cf)
	// write file back
	var encoder = yaml.NewEncoder(f)
	e = encoder.Encode(m)
	if e != nil {
		logger.Errorf("write config: %v", e)
		return e
	}

	return nil
}

func readRepoConfigFromPipe(ctx context.Context, reader io.Reader) (*configFileT, error) {
	var (
		logger  = logctx.Use(ctx).Sugar()
		decoder = yaml.NewDecoder(reader)
		e       error
		con     configFileT
	)
	e = decoder.Decode(&con)
	if errors.Is(e, io.EOF) {
		logger.Debugf("read config EOF, the file is empty?")
		e = nil
	}
	if e != nil {
		logger.Errorf("read config: %v", e)
		return nil, e
	}
	if con.Version != "" && con.Version != RepoConfigVersion {
		logger.Warnf("current config version: %v, not supported", con.Version)
		return nil, ErrRepoConfigUnsupportedVersion
	}
	return &con, nil
}

func repoConfigPath(repoPath string) string {
	return filepath.Join(repoPath, DefaultRepoConfigName)
}
