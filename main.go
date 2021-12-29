package main

import (
	"murphysec-cli-simple/cmd"
	"murphysec-cli-simple/logger"
)

func main() {
	cmd.Execute()
	logger.CloseAndWait()
}
