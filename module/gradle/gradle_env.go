package gradle

import (
	"context"
	"errors"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"os/exec"
)

//goland:noinspection GoNameStartsWithPackageName
type GradleEnv struct {
	Version             GradleVersion       `json:"version"`
	Path                string              `json:"path,omitempty"`
	GradleWrapperStatus GradleWrapperStatus `json:"gradle_wrapper_status"`
	GradleWrapperError  error               `json:"gradle_wrapper_error,omitempty"`
}

func (g *GradleEnv) ExecuteContext(ctx context.Context, args ...string) *exec.Cmd {
	var _args = make([]string, 0, len(args)+8)
	_args = append(_args, "--quiet", "--console", "plain")
	_args = append(_args, args...)
	c := exec.CommandContext(ctx, g.Path, _args...)
	logctx.Use(ctx).Sugar().Infof("Prepare: %s", c.String())
	return c
}

func DetectGradleEnv(ctx context.Context, dir string) (*GradleEnv, error) {
	var log = logctx.Use(ctx).Sugar()
	var r = &GradleEnv{GradleWrapperStatus: GradleWrapperStatusNotDetected}
	if script := prepareGradleWrapperScriptFile(ctx, dir); script != "" {
		gv, e := evalVersion(ctx, script)
		if e == nil {
			return &GradleEnv{
				Version:             *gv,
				Path:                script,
				GradleWrapperStatus: GradleWrapperStatusUsed,
			}, nil
		}
		log.Errorf("Eval gradle wrapper: %s", e.Error())
		r.GradleWrapperError = e
		r.GradleWrapperStatus = GradleWrapperStatusError
	}
	gv, e := evalVersion(ctx, "gradle")
	if e != nil {
		log.Errorf("Eval gradle: %s", e.Error())
		return nil, e
	}
	r.Version = *gv
	r.Path = "gradle"
	return r, nil
}

func evalVersion(ctx context.Context, cmdPath string) (_ *GradleVersion, err error) {
	defer func() {
		err = evalVersionError(err)
	}()
	var log = logctx.Use(ctx).Sugar()
	cmd := exec.CommandContext(ctx, cmdPath, "--version", "--quiet")
	log.Infof("Execute: %s", cmd.String())
	data, e := cmd.Output()
	if e != nil {
		var exitErr *exec.ExitError
		if errors.As(e, &exitErr) {
			data := exitErr.Stderr
			if len(data) > 256 {
				data = data[:256]
			}
			return nil, &EvalVersionError{
				_Error:   e,
				ExitCode: exitErr.ExitCode(),
				Stderr:   string(data),
			}
		}
		return nil, e
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
		if len(data) > 256 {
			data = data[:256]
		}
		return &EvalVersionError{
			_Error:   e,
			ExitCode: exitErr.ExitCode(),
			Stderr:   string(data),
		}
	}
	return &EvalVersionError{_Error: e}
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
	return fmt.Sprintf("%s, output: \n%s", e._Error.Error(), e.Stderr)
}

func (e *EvalVersionError) Is(target error) bool {
	return target == ErrEvalGradleVersion
}
