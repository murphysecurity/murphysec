package main

import (
	"github.com/ztrue/shutdown"
	"murphysec-cli-simple/cmd"
	"murphysec-cli-simple/logger"
	"os"
	"syscall"
)

func main() {
	go func() {
		shutdown.Listen(os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
		logger.Warn.Println("Signal received")
		os.Exit(-1)
	}()
	e := cmd.Execute()
	if e != nil {
		os.Exit(-1)
	}
	os.Exit(cmd.GetGlobalExitCode())
}
