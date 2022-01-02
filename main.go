package main

import (
	"fmt"
	"murphysec-cli-simple/cmd"
	"murphysec-cli-simple/util/output"
	"os"
	"strings"
)

func main() {
	output.Debug(fmt.Sprintf("CLI arguments: %s", strings.Join(os.Args, " ")))
	r := cmd.RootCmd()
	_ = r.Execute()
}
