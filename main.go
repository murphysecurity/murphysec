package main

import (
	"fmt"
	"github.com/ztrue/shutdown"
	"murphysec-cli-simple/cmd"
	"murphysec-cli-simple/util/output"
	"os"
	"strings"
)

func main() {
	go func() {
		shutdown.Listen(os.Interrupt, os.Kill)
		output.Error("User request interrupt.")
		os.Exit(-2)
	}()
	output.Debug(fmt.Sprintf("CLI arguments: %s", strings.Join(os.Args, " ")))
	r := cmd.RootCmd()
	_ = r.Execute()
}
