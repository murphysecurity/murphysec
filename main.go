package main

import (
	"github.com/murphysecurity/murphysec/cmd"
	"github.com/murphysecurity/murphysec/logger"
	"github.com/ztrue/shutdown"
	"os"
	"syscall"
)

func main() {
	go func() {
		shutdown.Listen(os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
		logger.Warn.Println("Signal received")
		os.Exit(-1)
	}()
	logger.LogFileCleanup()
	e := cmd.Execute()
	if e != nil {
		os.Exit(-1)
	}
	os.Exit(cmd.GetGlobalExitCode())
}
