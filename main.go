package main

import (
	"murphysec-cli-simple/cmd"
	"murphysec-cli-simple/logger"
	"os"
)

func main() {
	e := cmd.Execute()
	if e != nil {
		logger.Err.Println(e.Error())
	}
	logger.CloseAndWait()
	if e != nil {
		os.Exit(-1)
	}
	os.Exit(cmd.GetGlobalExitCode())
}
