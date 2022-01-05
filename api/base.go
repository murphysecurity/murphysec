package api

import (
	"fmt"
	"murphysec-cli-simple/util/output"
	"os"
	"strings"
)

// serverAddress returns API base URL
func serverAddress() string {
	envServer := strings.Trim(strings.TrimSpace(os.Getenv("MPS_CLI_SERVER")), "/")
	if len(envServer) == 0 {
		return "https://www.murphysec.com/api"
	}
	return envServer
}
func init() {
	output.Debug(fmt.Sprintf("Server addr: %s", serverAddress()))
}

var defaultToken string

func SetDefaultToken(token string) {
	defaultToken = token
}
