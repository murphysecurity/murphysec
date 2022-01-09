//go:build !windows
// +build !windows

package utils

import (
	"murphysec-cli-simple/logger"
	"os/exec"
	"strconv"
)

func KillAllChild(ppid int) {
	if ppid < 0 {
		return
	}
	c := exec.Command("pkill", "-15", "-p", strconv.Itoa(ppid))
	logger.Debug.Printf("execute: %s", c.String())
	_, _ = c.Output()
	logger.Debug.Printf("Done")
}
