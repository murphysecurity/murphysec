package util

import (
	"io/ioutil"
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
	}()

	go func() {
		all, err := ioutil.ReadAll(stde)
		t.stderr.ch <- streamR{t: string(all), e: err}
		close(t.stderr.ch)
	}()
	if e := t.cmd.Start(); e != nil {
		return e
	}
	t.pid = t.cmd.Process.Pid

	go func() {
		select {
		case <-t.abortChan:
			p := t.cmd.Process
			if p != nil {
				_ = p.Kill()
			}
		case <-finishChan:
		}
	}()
	err := t.cmd.Wait()
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
