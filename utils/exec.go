package utils

import (
	"github.com/pkg/errors"
	"github.com/ztrue/shutdown"
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"os"
	"os/exec"
	"sync"
)

func ExecuteCmd(cmd string, arg ...string) *PreparedCmd {
	return &PreparedCmd{
		cmd:       exec.Command(cmd, arg...),
		abortChan: make(chan struct{}),
	}
}

type PreparedCmd struct {
	abortChan  chan struct{}
	cmd        *exec.Cmd
	pid        int
	stdoutData []byte
	stderrData []byte
	stdoutErr  error
	stderrErr  error
}

func (this *PreparedCmd) Abort() {
	p := this.cmd.Process
	if p != nil {
		logger.Debug.Printf("Abort command[pid=%d]: %s", p.Pid, this.cmd.String())
	} else {
		logger.Debug.Printf("Abort command: %s", this.cmd.String())
	}
	close(this.abortChan)
}

func (this *PreparedCmd) Execute() error {
	logger.Debug.Printf("Execute cmd: %s", this.cmd.String())
	//get stdout & stderr
	stdout, e := this.cmd.StdoutPipe()
	if e != nil {
		return errors.Wrap(e, "Get stdout failed")
	}
	stderr, e := this.cmd.StderrPipe()
	if e != nil {
		return errors.Wrap(e, "Get stderr failed")
	}
	// executing
	if e := this.cmd.Start(); e != nil {
		return errors.Wrap(e, "Execute command failed")
	}
	shutdownKey := shutdown.Add(func() {
		pid := this.pid
		process := this.cmd.Process
		if pid == 0 || process == nil {
			return
		}
		logger.Debug.Printf("Sending interrupt signal, pid: %d", pid)
		if e := this.cmd.Process.Signal(os.Interrupt); e != nil {
			logger.Debug.Printf("Sending interrupt signal failed, use kill. %s", e.Error())
			if e := process.Kill(); e != nil {
				logger.Debug.Printf("Kill failed. %s", e.Error())
			}
			KillAllChild(pid)
		}
	})
	defer shutdown.Remove(shutdownKey)
	this.pid = this.cmd.Process.Pid
	logger.Debug.Printf("Command started, pid: %d", this.Pid())
	// collect output
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		this.stdoutData, this.stdoutErr = ioutil.ReadAll(stdout)
		if this.stdoutErr != nil {
			logger.Debug.Printf("Pid: %d -> Stdout read, err: %s", this.Pid(), this.stdoutErr.Error())
		} else {
			logger.Debug.Printf("Pid: %d -> Stdout read with no errors", this.Pid())
		}
		wg.Done()
	}()
	go func() {
		this.stderrData, this.stderrErr = ioutil.ReadAll(stderr)
		if this.stderrErr != nil {
			logger.Debug.Printf("Pid: %d -> Stderr read, err: %s", this.Pid(), this.stderrErr.Error())
		} else {
			logger.Debug.Printf("Pid: %d -> Stderr read with no errors", this.Pid())
		}
		wg.Done()
	}()
	wg.Wait()
	if e := this.cmd.Wait(); e == nil {
		logger.Debug.Printf("Pid: %d -> Execution terminated with no errors.", this.Pid())
	} else {
		logger.Debug.Printf("Pid: %d -> Execution terminated with err: %s", this.Pid(), e.Error())
		return e
	}
	return nil
}

func (this *PreparedCmd) GetStdout() (string, error) {
	if this.stdoutErr != nil {
		return "", this.stdoutErr
	}
	return string(this.stdoutData), nil
}

func (this *PreparedCmd) GetStderr() (string, error) {
	if this.stderrErr != nil {
		return "", this.stderrErr
	}
	return string(this.stderrData), nil
}

func (this *PreparedCmd) Pid() int {
	return this.pid
}
