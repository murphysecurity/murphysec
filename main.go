package main

import (
	"github.com/murphysecurity/murphysec/cmd"
	"os"
)

func main() {
	e := cmd.Execute()
	if e != nil {
		os.Exit(-1)
	}
	os.Exit(cmd.GetGlobalExitCode())
}
