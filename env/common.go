package env

import (
	"os"
	"strconv"
)

func envi(name string, defaultValue int) int {
	if i, e := strconv.Atoi(os.Getenv(name)); e != nil {
		return defaultValue
	} else {
		return i
	}
}

const DefaultServerURL = "https://www.murphysec.com"

var ServerURLOverride = os.Getenv("MPS_CLI_SERVER")
var APITokenOverride = os.Getenv("API_TOKEN")
