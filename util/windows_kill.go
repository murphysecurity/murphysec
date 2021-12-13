//go:build windows
// +build windows

package util

import (
	"fmt"
	"murphysec-cli-simple/util/output"
	"os/exec"
	"strconv"
)

func KillAllChild(ppid int) {
	c := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(ppid))
	output.Debug(fmt.Sprintf("execute: %s", c.String()))
	_, _ = c.Output()
	output.Debug("Done")
}
