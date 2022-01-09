package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func WatchKill() (chan bool, func()) {
	rsCh := make(chan bool, 1)
	termCh := make(chan os.Signal, 1)
	finishCh := make(chan struct{})
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	canceller := func() {
		signal.Stop(termCh)
		close(finishCh)
	}
	go func() {
		select {
		case <-termCh:
			rsCh <- true
		case <-finishCh:
			rsCh <- false
		}
	}()
	return rsCh, canceller
}
