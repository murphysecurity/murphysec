package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os/exec"
	"strings"
	"time"
)

var ErrMvnExitErr = mvnError("mvn command exit with non-zero code")
var ErrMvnCmd = mvnError("error during mvn execution")

type MvnGraphCmdArgs struct {
	Path     string
	Profiles []string
	Timeout  time.Duration
	ScanDir  string
}

func (m MvnGraphCmdArgs) Execute(ctx context.Context) error {
	logger := utils.UseLogger(ctx)
	if ctx == nil {
		ctx = context.TODO()
	}
	if m.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, m.Timeout)
		defer cancel()
	}
	var args = []string{"com.github.ferstl:depgraph-maven-plugin:4.0.1:graph", "-DgraphFormat=json", "--batch-mode"}
	if len(m.Profiles) > 0 {
		args = append(args, "-P")
		args = append(args, strings.Join(m.Profiles, ","))
	}
	c := exec.CommandContext(ctx, m.Path, args...)
	c.Dir = m.ScanDir
	logStream := utils.NewLogPipe(logger, "mvn")
	defer logStream.Close()
	c.Stderr = logStream
	c.Stdout = logStream
	logger.Info(fmt.Sprintf("Start command: %s", c.String()), zap.String("dir", c.Dir))
	if e := c.Start(); e != nil {
		logger.Error("Start command failed", zap.Error(e))
		return errors.WithMessage(ErrMvnCmd, e.Error())
	}
	if e := c.Wait(); e != nil {
		logger.Error(ErrMvnCmd.Error(), zap.Error(e), zap.Int("exit_code", c.ProcessState.ExitCode()))
		return errors.WithMessage(ErrMvnCmd, e.Error())
	}
	exitCode := c.ProcessState.ExitCode()
	if exitCode != 0 {
		logger.Error(ErrMvnExitErr.Error(), zap.Int("code", exitCode))
		return errors.WithMessage(ErrMvnExitErr, fmt.Sprintf("code: %d", exitCode))
	}
	logger.Info("Mvn graph command exit with no errors")
	return nil
}
