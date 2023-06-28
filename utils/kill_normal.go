//go:build !windows

package utils

import (
	"os/exec"
	"syscall"
)

func KillProcessGroup(pid int) {
	_ = syscall.Kill(-pid, syscall.SIGKILL)
}

func SetPGid(c *exec.Cmd) {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Setpgid = true
}
