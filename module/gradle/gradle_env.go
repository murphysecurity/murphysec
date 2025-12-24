package gradle

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/module/gradle/bundle"
)

//goland:noinspection GoNameStartsWithPackageName
type GradleEnv struct {
	Version             GradleVersion       `json:"version"`
	Path                string              `json:"path,omitempty"`
	GradleWrapperStatus GradleWrapperStatus `json:"gradle_wrapper_status"`
	GradleWrapperError  error               `json:"gradle_wrapper_error,omitempty"`
	JavaHome            string              `json:"java_home,omitempty"`
}

func (g *GradleEnv) ExecuteContext(ctx context.Context, args ...string) *exec.Cmd {
	var _args = make([]string, 0, len(args)+8)
	_args = append(_args, "--info", "--console", "plain", "--stacktrace")
	_args = append(_args, args...)
	c := exec.CommandContext(ctx, g.Path, _args...)
	c.Env = os.Environ()
	if g.JavaHome != "" {
		c.Env = append(c.Env, "JAVA_HOME="+g.JavaHome)
		c.Env = append(c.Env, "PATH="+filepath.Join(g.JavaHome, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
	}
	logctx.Use(ctx).Sugar().Infof("Prepare: %s", c.String())
	return c
}

func DetectGradleEnv(ctx context.Context, dir string) (*GradleEnv, error) {
	var log = logctx.Use(ctx).Sugar()
	var r = &GradleEnv{GradleWrapperStatus: GradleWrapperStatusNotDetected}
	var gwv = readGradleVersionFromWrapper(ctx, dir)
	if os.Getenv("MPS_BUNDLED_GRADLE") == "1" {
		if gwv == "" {
			// no version read, use default latest
			log.Info("use default gradle version")
			gwv = "8.6"
		}
		r.Path, r.JavaHome = bundle.FindOkVersion(ctx, gwv)
		log.Infof("use bundled gradle: %v", r.Path)
		log.Infof("use bundled java: %v", r.JavaHome)
	} else if os.Getenv("MPS_BUNDLED_JAVA") == "1" {
		r.JavaHome = bundle.SelectJavaHome(gwv)
		log.Infof("use bundled Java eval gradle version(from wrapper): %s", r.JavaHome)
	}
	if r.Path != "" {
		return r, nil
	}
	if script := prepareGradleWrapperScriptFile(ctx, dir); script != "" {
		gv, e := evalVersion(ctx, script, r.JavaHome)
		if e == nil {
			r.Version = *gv
			r.Path = script
			r.GradleWrapperStatus = GradleWrapperStatusUsed
			return r, nil
		}
		log.Errorf("Eval gradle wrapper: %s", e.Error())
		r.GradleWrapperError = e
		r.GradleWrapperStatus = GradleWrapperStatusError
	}
	gv, e := evalVersion(ctx, "gradle", "")
	if e != nil {
		log.Errorf("Eval gradle: %s", e.Error())
		return nil, e
	}
	r.Version = *gv
	r.Path = "gradle"
	return r, nil
}

func evalVersion(ctx context.Context, cmdPath string, javaHome string) (_ *GradleVersion, err error) {
	defer func() {
		err = evalVersionError(err)
	}()
	var log = logctx.Use(ctx).Sugar()
	cmd := exec.CommandContext(ctx, cmdPath, "--version", "--quiet")
	if javaHome != "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "JAVA_HOME="+javaHome)
	}
	log.Infof("Execute: %s", cmd.String())
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseGradleVersionOutput(string(data))
}

func evalVersionError(e error) error {
	if e == nil {
		return nil
	}
	var exitErr *exec.ExitError
	if errors.As(e, &exitErr) {
		data := exitErr.Stderr
		var s = strings.TrimSpace(string(data))
		if len(s) == 0 {
			s = "(no stderr output)"
		} else if len(s) > 1024 {
			s = s[:1024]
		}
		return &EvalVersionError{
			_Error:   e,
			ExitCode: exitErr.ExitCode(),
			Stderr:   s,
		}
	}
	return &EvalVersionError{_Error: e}
}

var gradlePropertiesWrapperVer *regexp.Regexp
var gradlePropertiesWrapperVerOnce sync.Once

func readGradleVersionFromWrapper(ctx context.Context, projectDir string) string {
	gradlePropertiesWrapperVerOnce.Do(func() {
		gradlePropertiesWrapperVer = regexp.MustCompile(`distributionUrl=.+?gradle-([0-9.]+)(?:-bin|-all|).zip`)
	})
	var logger = logctx.Use(ctx).Sugar()
	var propertiesFilepath = filepath.Join(projectDir, "gradle", "wrapper", "gradle-wrapper.properties")
	data, e := os.ReadFile(propertiesFilepath)
	if e != nil {
		if os.IsNotExist(e) {
			logger.Debug("gradle-wrapper.properties doesn't exists")
		} else {
			logger.Warnf("read gradle-wrapper.properties failed: %e", e)
		}
		return ""
	}
	var r = gradlePropertiesWrapperVer.FindStringSubmatch(string(data))
	if len(r) > 0 {
		logger.Debugf("gradle-wrapper gradle version read: %s", r[1])
		return r[1]
	}
	return ""
}

type EvalVersionError struct {
	_Error   error
	ExitCode int    `json:"exit_code"`
	Stderr   string `json:"stderr"`
}

func (e *EvalVersionError) Unwrap() error {
	return e._Error
}

func (e *EvalVersionError) Error() string {
	if e.Stderr == "" {
		return e._Error.Error()
	}
	return fmt.Sprintf("%s, output: \n%s", e._Error.Error(), strconv.Quote(e.Stderr))
}

func (e *EvalVersionError) Is(target error) bool {
	return target == ErrEvalGradleVersion
}
