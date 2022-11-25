package main

import (
	"fmt"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"os"
)

func main() {
	if e := rootCmd().Execute(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
	exitcode.Exit()
}
