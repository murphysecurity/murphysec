package envinspection

import (
	"bytes"
	"context"
	"os/exec"
)

type cError struct {
	Content string
	Err     error
}

func (e cError) Error() string {
	return e.Err.Error()
}

func handleCmd(ctx context.Context, cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", cError{
			Content: stderr.String(),
			Err:     err,
		}
	}
	return stdout.String(), nil
}
