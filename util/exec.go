package util

import (
	"fmt"
	"io/ioutil"
	"murphysec-cli-simple/util/output"
	"os/exec"
	"sync"
)

func ExecuteCmd(cmd string, arg ...string) *PreparedCmd {
	return &PreparedCmd{
		cmd:       exec.Command(cmd, arg...),
		abortChan: make(chan struct{}),
		stderr:    newStream(),
		stdout:    newStream(),
	}
}

type PreparedCmd struct {
	abortChan chan struct{}
	cmd       *exec.Cmd
	stdout    *stream
	stderr    *stream
	pid       int
}

type stream struct {
	text string
	err  error
	once sync.Once
	ch   chan streamR
}

func newStream() *stream {
	return &stream{
		text: "",
		err:  nil,
		once: sync.Once{},
		ch:   make(chan streamR),
	}
}

func (s *stream) read() (string, error) {
	s.once.Do(func() {
		c := <-s.ch
		s.text = c.t
		s.err = c.e
	})
	return s.text, s.err
}

type streamR struct {
	t string
	e error
}

func (t *PreparedCmd) Abort() {
	p := t.cmd.Process
	if p != nil {
		output.Debug(fmt.Sprintf("Abort command[pid=%d]: %s", p.Pid, t.cmd.String()))
	} else {
		output.Debug(fmt.Sprintf("Abort command: %s", t.cmd.String()))
	}
	close(t.abortChan)
}

func (t *PreparedCmd) Execute() error {
	finishChan := make(chan struct{})
	stdo, e := t.cmd.StdoutPipe()
	if e != nil {
		return e
	}
	stde, e := t.cmd.StderrPipe()
	if e != nil {
		return e
	}

	go func() {
		all, err := ioutil.ReadAll(stdo)
		t.stdout.ch <- streamR{t: string(all), e: err}
		close(t.stdout.ch)
		output.Debug(fmt.Sprintf("Process %d stdout read", t.pid))
	}()

	go func() {
		all, err := ioutil.ReadAll(stde)
		t.stderr.ch <- streamR{t: string(all), e: err}
		close(t.stderr.ch)
		output.Debug(fmt.Sprintf("Process %d stderr read", t.pid))
	}()
	if e := t.cmd.Start(); e != nil {
		return e
	}
	t.pid = t.cmd.Process.Pid
	output.Debug(fmt.Sprintf("Cmd %s started, pid: %d", t.cmd.String(), t.pid))
	go func() {
		select {
		case <-t.abortChan:
			p := t.cmd.Process
			if p != nil {
				output.Debug(fmt.Sprintf("Kill process: %d", p.Pid))
				if e := p.Kill(); e == nil {
					output.Debug("Kill succeed")
				} else {
					output.Warn(fmt.Sprintf("Kill process failed: %s", e.Error()))
				}
			}
		case <-finishChan:
			output.Debug(fmt.Sprintf("Process %d finished", t.pid))
		}
	}()
	err := t.cmd.Wait()
	if err == nil {
		output.Debug(fmt.Sprintf("Process %d wait finished with no err", t.pid))
	} else {
		output.Debug(fmt.Sprintf("Process %d wait finished with err: %s", t.pid, err.Error()))
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *PreparedCmd) GetStdout() (string, error) {
	return t.stdout.read()
}

func (t *PreparedCmd) GetStderr() (string, error) {
	return t.stderr.read()
}

func (t *PreparedCmd) Pid() int {
	return t.pid
}
