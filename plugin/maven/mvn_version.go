package maven

import (
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/util"
	"murphysec-cli-simple/util/output"
	"strings"
)

func mavenVersion() (*RuntimeMavenVersion, error) {
	c := util.ExecuteCmd("mvn", "--version")
	killSig, canceller := util.WatchKill()
	defer canceller()
	go func() {
		if <-killSig {
			util.KillAllChild(c.Pid())
			c.Abort()
		}
	}()
	if e := c.Execute(); e != nil {
		fmt.Println(e.Error())
		if s, e := c.GetStderr(); e != nil {
			output.Warn(fmt.Sprintf("Get error out failed: %s", e.Error()))
		} else {
			output.Warn(s)
		}
		return nil, errors.Wrap(e, "Get maven version failed")
	}
	if t, e := c.GetStdout(); e == nil {
		return parseMvnVerCommandResult(t), nil
	} else {
		return nil, errors.Wrap(e, "Read maven stdout failed")
	}
}

type RuntimeMavenVersion struct {
	MvnVersion  string `json:"mvn_version"`
	JavaVersion string `json:"java_version"`
	RuntimeOs   string `json:"runtime_os"`
}

func parseMvnVerCommandResult(cmdResult string) *RuntimeMavenVersion {
	lines := strings.Split(cmdResult, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	rs := RuntimeMavenVersion{}
	for _, it := range lines {
		switch {
		case strings.HasPrefix(it, "Apache Maven"):
			rs.MvnVersion = it
		case strings.HasPrefix(it, "Java version"):
			rs.JavaVersion = it
		case strings.HasPrefix(it, "Os name"):
			rs.RuntimeOs = it
		}
	}
	return &rs
}
