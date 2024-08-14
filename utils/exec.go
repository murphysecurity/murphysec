package utils

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

func ExecGetStdOutErr(ctx context.Context, cmd *exec.Cmd) (stdout io.ReadCloser, stderr io.ReadCloser, e error) {
	stdout, e = cmd.StdoutPipe()
	if e != nil {
		e = fmt.Errorf("cmd.StdoutPipe(): %w", e)
		return
	}
	stderr, e = cmd.StderrPipe()
	if e != nil {
		e = fmt.Errorf("cmd.StderrPipe(): %w", e)
		return
	}
	return
}
