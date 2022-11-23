package main

import (
	"github.com/murphysecurity/murphysec/cmd"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"os"
)

func main() {
	e := cmd.Execute()
	if e != nil {
		os.Exit(-1)
	}
	exitcode.Exit()
}
