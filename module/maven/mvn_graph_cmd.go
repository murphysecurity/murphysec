package maven

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/infra/logpipe"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

// PluginGraphCmd helper to com.github.ferstl:depgraph-maven-plugin:4.0.1:graph
type PluginGraphCmd struct {
	Profiles     []string
	Timeout      time.Duration
	ScanDir      string
	MavenCmdInfo *MvnCommandInfo
}

func (m PluginGraphCmd) RunC(ctx context.Context) error {
	logger := logctx.Use(ctx)
	if ctx == nil {
		ctx = context.TODO()
	}
	if m.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, m.Timeout)
		defer cancel()
	}
	var args = []string{"com.github.ferstl:depgraph-maven-plugin:4.0.1:graph", "-DgraphFormat=json"}
	if len(m.Profiles) > 0 {
		args = append(args, "-P")
		args = append(args, strings.Join(m.Profiles, ","))
	}
	c := m.MavenCmdInfo.Command(ctx, args...)
	c.Dir = m.ScanDir
	logStream := logpipe.New(logger, "mvn")
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
