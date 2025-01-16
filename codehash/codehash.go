package codehash

import (
	"bufio"
	"context"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/must"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Scan(ctx context.Context) {
	var logger = logctx.Use(ctx)
	var base = filepath.Dir(must.A(os.Executable()))

	var st = model.UseScanTask(ctx)
	tempFile := must.A(os.CreateTemp("", "codehash-output-*.json"))
	var tempFilePath = tempFile.Name()
	_ = tempFile.Close()
	_ = os.Remove(tempFilePath)
	defer func() { _ = os.Remove(tempFilePath) }()
	var cmd *exec.Cmd
	var shPath string
	if runtime.GOOS == "windows" {
		shPath = filepath.Join(base, "bin\\cli.bat")
		cmd = exec.CommandContext(ctx, shPath, st.ProjectPath, tempFilePath, st.SubtaskId)
	} else {
		shPath = filepath.Join(base, "bin/cli")
		fi := must.A(os.Stat(shPath))
		must.Must(os.Chmod(shPath, fi.Mode()|0111))
		cmd = exec.CommandContext(ctx, "sh", filepath.Join(base, "bin/cli"), st.ProjectPath, tempFilePath, st.SubtaskId)
	}
	stderr := must.A(cmd.StderrPipe())
	go func() {
		defer func() { _ = stderr.Close() }()
		b := bufio.NewScanner(stderr)
		for b.Scan() {
			logger.Sugar().Infof("codehash[e]: %s", b.Text())
		}
	}()
	e := cmd.Run()
	if e != nil {
		logger.Sugar().Errorf("codehash error: %s", e.Error())
		return
	}
	e = api.SubmitCodeHash(ctx, api.DefaultClient(), tempFilePath)
	if e != nil {
		logger.Sugar().Errorf("codehash submit error: %s", e.Error())
		return
	}
}
