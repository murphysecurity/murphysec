package config

//go:generate stringer -type repoConfigError -linecomment -output repo_config_error_string.go
type repoConfigError int

const (
	_                               repoConfigError = iota
	ErrRepoConfigNotFound                           // config: no repo config found
	ErrRepoConfigBad                                // config: repo config broken
	ErrRepoConfigUnsupportedVersion                 // config: repo config version unsupported
	ErrRepoConfigBadAccessType                      // config: bad access type
)

func (i repoConfigError) Error() string {
	return i.String()
}
