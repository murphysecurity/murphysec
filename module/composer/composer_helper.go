package composer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var ErrComposerVersionCheckFail = errors.New("Check composer version failed")
var ErrComposerResolveFail = errors.New("PHP composer resolve failed")

func doComposerInstall(ctx context.Context, projectDir string) error {
	if ctx == nil {
		ctx = context.TODO()
	}
	c := exec.CommandContext(ctx, "composer", "--ignore-platform-reqs", "--no-progress", "--no-dev", "--no-autoloader", "--no-scripts", "--no-interaction", "--quiet")
	c.Dir = projectDir
	logger.Info.Println("Command:", c.String())
	cif := &composerInstallFail{}
	c.Stderr = cif
	if e := c.Run(); e == nil {
		return nil
	} else {
		cif.err = e
	}
	exitCode := c.ProcessState.ExitCode()
	cif.code = exitCode

	if exitCode == 2 {
		return ErrComposerResolveFail
	}
	return cif
}

type composerInstallFail struct {
	code         int
	stdErrPrefix []byte
	err          error
}

func (c *composerInstallFail) Write(data []byte) (int, error) {
	const maxPrefix = 1024
	c.stdErrPrefix = append(c.stdErrPrefix, data[:maxPrefix-len(c.stdErrPrefix)]...)
	return len(data), nil
}

func (c composerInstallFail) Unwrap() error {
	return c.err
}

func (c composerInstallFail) requirementsCouldNotBeResolved() bool {
	return strings.Contains(string(c.stdErrPrefix), "Your requirements could not be resolved")
}

func (c composerInstallFail) Error() string {
	if c.requirementsCouldNotBeResolved() {
		return ErrComposerResolveFail.Error()
	}
	if c.err == nil {
		return fmt.Sprintf("Composer exit with error: %s", strconv.Quote(string(c.stdErrPrefix)))
	}
	return fmt.Sprintf("Composer: %s %s", c.err.Error(), strconv.Quote(string(c.stdErrPrefix)))
}

func (c composerInstallFail) Is(target error) bool {
	return target == ErrComposerResolveFail
}

var composerVersionPattern = regexp.MustCompile("Composer version ([^ ]+)")

func checkComposerVersion() (interface{}, error) {
	c := exec.Command("composer", "--version")
	logger.Info.Println("Command:", c.String())
	data, e := c.Output()
	if e != nil {
		return nil, errors.WithStack(composerVersionCheckFail(e.Error()))
	}
	m := composerVersionPattern.FindStringSubmatch(string(data))
	if m == nil {
		return nil, errors.WithStack(composerVersionCheckFail("Version pattern not match"))
	}
	return &PhpComposerVersion{Version: m[1]}, nil
}

type PhpComposerVersion struct {
	Version string
}

func (receiver PhpComposerVersion) String() string {
	return receiver.Version
}

type composerVersionCheckFail string

func (r composerVersionCheckFail) Error() string {
	return fmt.Sprintf("ComposerVersionCheckFail: %s", string(r))
}

func (r composerVersionCheckFail) Is(target error) bool {
	return r == ErrComposerVersionCheckFail || r == target
}
