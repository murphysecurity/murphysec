package main

import (
	"github.com/murphysecurity/murphysec/cmd"
	"github.com/murphysecurity/murphysec/logger"
	"os"
)

func main() {
	logger.LogFileCleanup()
	e := cmd.Execute()
	if e != nil {
		os.Exit(-1)
	}
	os.Exit(cmd.GetGlobalExitCode())
}
