package env

import (
	"os"
	"strconv"
)

var GradleExecutionTimeoutSecond = envi("GRADLE_EXECUTION_TIMEOUT_SEC", 20*60)

func envi(name string, defaultValue int) int {
	if i, e := strconv.Atoi(os.Getenv(name)); e != nil {
		return defaultValue
	} else {
		return i
	}
}
