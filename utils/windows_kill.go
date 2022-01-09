//go:build windows
// +build windows

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
	c := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(ppid))
	logger.Debug.Printf("execute: %s", c.String())
	_, _ = c.Output()
	logger.Debug.Printf("Done")
}
