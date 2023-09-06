package python

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"go.uber.org/zap"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

var __oncePipList sync.Once
var __oncePipListM map[string]string
var __oncePipListE error

func getEnvPipListMap(ctx context.Context) (map[string]string, error) {
	__oncePipList.Do(func() {
		__oncePipListM, __oncePipListE = evalPipList(ctx)
	})
	return __oncePipListM, __oncePipListE
}

func evalPipList(ctx context.Context) (map[string]string, error) {
	const defaultPipTimeout = time.Second * 15
	pip, e := selectPip(ctx)
	if e != nil {
		return nil, e
	}
	m, e := evalPipListJsonFormat(ctx, pip, defaultPipTimeout)
	if e == nil {
		return m, nil
	}
	m, e = evalPipListDefaultFormat(ctx, pip, defaultPipTimeout)
	return m, e
}

func evalPipListDefaultFormat(ctx context.Context, pip string, timeout time.Duration) (map[string]string, error) {
	data, e := executePipListCmd(ctx, timeout, pip, "")
	if e != nil {
		return nil, e
	}
	return parsePipListDefaultFormat(data)
}

func evalPipListJsonFormat(ctx context.Context, pip string, timeout time.Duration) (map[string]string, error) {
	data, e := executePipListCmd(ctx, timeout, pip, "json")
	if e != nil {
		return nil, e
	}
	return parsePipListJson(data)
}

func selectPip(ctx context.Context) (string, error) {
	if env.DoNotBuild {
		return "", fmt.Errorf("pip: do not build")
	}
	if checkCmdExists(ctx, "pip3") {
		return "pip3", nil
	}
	if checkCmdExists(ctx, "pip") {
		return "pip", nil
	}
	return "", fmt.Errorf("no pip command")
}

var __pipListDefaultPattern = regexp.MustCompile(`^([\w_-]+)\s\(([\w.-]+)\)$`)

func parsePipListDefaultFormat(data []byte) (map[string]string, error) {
	var r = make(map[string]string)
	for _, s := range strings.Split(string(data), "\n") {
		m := __pipListDefaultPattern.FindStringSubmatch(strings.TrimSpace(s))
		if len(m) == 0 {
			continue
		}
		r[m[1]] = m[2]
	}
	return r, nil
}

func parsePipListJson(data []byte) (map[string]string, error) {
	type t struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	var r []t
	if e := json.Unmarshal(data, &r); e != nil {
		return nil, e
	}
	sort.Slice(r, func(i, j int) bool {
		if r[i].Name == r[j].Name {
			return r[i].Version < r[j].Version
		}
		return r[i].Name < r[j].Name
	})
	var m = make(map[string]string)
	for _, it := range r {
		if it.Version == "" {
			continue
		}
		m[it.Name] = it.Version
	}
	return m, nil
}

func executePipListCmd(ctx context.Context, timeout time.Duration, pipCmd string, format string) ([]byte, error) {
	var logger = logctx.Use(ctx).Sugar()
	var args = []string{"list", "--disable-pip-version-check", "--no-index"}
	if format != "" {
		args = append(args, "--format", format)
	}
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(cctx, pipCmd, args...)
	logger.Infof("eval piplist: %v", cmd)
	output, e := cmd.Output()
	if e != nil {
		var ee *exec.ExitError
		if errors.As(e, &ee) {
			printExitErrorLog(logger, ee)
			return nil, fmt.Errorf("eval pip list failed")
		}
	}
	logger.Infof("execute succeeded, read %d bytes", len(output))
	return output, nil
}

func printExitErrorLog(logger *zap.SugaredLogger, e *exec.ExitError) {
	logger.Warnf("exit with error, code: %d", e.ExitCode())
	for _, s := range strings.Split(string(e.Stderr), "\n") {
		if strings.TrimSpace(s) == "" {
			continue
		}
		logger.Debugf("output: %s", s)
	}
}

func checkCmdExists(ctx context.Context, cmd string) bool {
	var logger = logctx.Use(ctx).Sugar()
	s, e := exec.LookPath(cmd)
	if e != nil || s == "" {
		logger.Infof("%s doesn't exists, %v", cmd, e)
		return false
	} else {
		logger.Infof("%s found: %s", cmd, s)
		return true
	}
}
