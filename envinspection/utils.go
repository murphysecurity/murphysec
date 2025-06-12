package envinspection

import (
	"bytes"
	"context"
	"os/exec"
)

func handleCmd(ctx context.Context, cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String(), err
}
