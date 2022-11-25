package buildinfo

import (
	"runtime/debug"
	"strconv"
	"time"
)

var version string

const devVersion = "(devel)"

type D struct {
	Version    string    `json:"version" yaml:"version"`
	CommitHash string    `json:"commit_hash" yaml:"commit_hash"`
	CommitTime time.Time `json:"commit_time" yaml:"commit_time"`
	Modified   bool      `json:"modified" yaml:"modified"`
}

func (d D) IsZero() bool {
	return d.Version == "" && d.CommitTime.IsZero() && d.CommitHash == "" && !d.Modified
}

func buildInfo() (d D) {
	d.Version = devVersion
	if version != "" {
		d.Version = version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	var m = map[string]string{}
	for _, setting := range info.Settings {
		m[setting.Key] = setting.Value
	}
	if m["vcs"] == "git" {
		d.Modified, _ = strconv.ParseBool(m["vcs.modified"])
		d.CommitTime, _ = time.Parse(time.RFC3339, m["vcs.time"])
		d.CommitHash = m["vcs.revision"]
	}
	return
}

var _d = buildInfo()

func Get() D {
	return _d // copied
}

// Commit returns a string contains git commit info
//
// example: 76f50a8c7ab9c3e8c488a597ec3fe173d6012dd5(modified) at 2006-01-02T15:04:05Z07:00
func Commit() string {
	d := Get()
	if d.CommitHash == "" {
		return ""
	}
	var s string
	if d.Modified {
		s = d.CommitHash + "(modified)"
	} else {
		s = d.CommitHash
	}
	if !d.CommitTime.IsZero() {
		s += " at " + d.CommitTime.Format(time.RFC3339)
	}
	return s
}

// UserAgentSuffix returns part of user-agent
//
// example: Commit/01234ef-modified
func UserAgentSuffix() (r string) {
	d := Get()
	if d.CommitHash != "" || len(d.CommitHash) != 40 {
		return
	}
	r += "Commit/"
	r += d.CommitHash[:7]
	if d.Modified {
		r += "-modified"
	}
	return
}
