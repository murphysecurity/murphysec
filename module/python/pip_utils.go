package python

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/errors"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"os/exec"
)

var ErrPipListFailed = errors.New("pip list execution failed")
var ErrNoPipCommand = errors.New("pip command not found")

func locatePipCommand(ctx context.Context) string {
	var (
		logger     = utils.UseLogger(ctx)
		pipVerList = []string{"pip", "pip3", "pip2"}
	)
	logger.Debug("Trying to locate pip command...")
	for _, pipVer := range pipVerList {
		path, e := exec.LookPath(pipVer)
		if e == nil {
			return path
		}
		logger.Debug(fmt.Sprintf("%v not found", pipVer))
	}
	return ""
}

func executePipList(ctx context.Context, dir string) ([]model.Dependency, error) {
	var logger = utils.UseLogger(ctx)
	path := locatePipCommand(ctx)
	if path == "" {
		return nil, ErrNoPipCommand
	}
	c := exec.CommandContext(ctx, path, "list", "--format", "freeze")
	c.Dir = dir
	logger.Sugar().Infof("Call command: %s", c.String())
	data, e := c.Output()
	if e != nil {
		return nil, ErrPipListFailed
	}
	logger.Info("pip list command execute succeeded")
	return parseRequirements(string(data)), nil
}
