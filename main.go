package main

import (
	"murphysec-cli-simple/cmd"
)

func main() {
	r := cmd.RootCmd()
	_ = r.Execute()
}
