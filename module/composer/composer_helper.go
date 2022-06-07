package composer

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/logger"
	"os/exec"
	"strconv"
)

const _ComposerInstallFailMaxPrefix = 1024

type composerInstallFail struct {
	exitCode     int
	stdErrPrefix []byte
	execError    error
}

func doComposerInstall(ctx context.Context, projectDir string) error {
	c := exec.CommandContext(ctx, "composer", "--ignore-platform-reqs", "--no-progress", "--no-dev", "--no-autoloader", "--no-scripts", "--no-interaction", "--quiet")
	c.Dir = projectDir
	logger.Info.Println("Command:", c.String())
	cif := &composerInstallFail{}
	c.Stderr = cif
	if e := c.Run(); e != nil {
		cif.execError = e
		cif.exitCode = c.ProcessState.ExitCode()
		if cif.exitCode == 2 {
			return wrapErr(ErrComposerResolveFail, cif)
		}
		return cif
	}
	return nil
}

func (c *composerInstallFail) Write(data []byte) (int, error) {
	const maxPrefix = _ComposerInstallFailMaxPrefix
	c.stdErrPrefix = append(c.stdErrPrefix, data[:maxPrefix-len(c.stdErrPrefix)]...)
	return len(data), nil
}

func (c composerInstallFail) Unwrap() error {
	return c.execError
}

func (c composerInstallFail) Error() string {
	if c.execError == nil {
		return fmt.Sprintf("Composer exit with error: %s", strconv.Quote(string(c.stdErrPrefix)))
	}
	return fmt.Sprintf("Composer: %s %s", c.execError.Error(), strconv.Quote(string(c.stdErrPrefix)))
}

func (c composerInstallFail) Is(target error) bool {
	return target == ErrComposerResolveFail
}
