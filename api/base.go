package api

import (
	"os"
	"strings"
)

// serverAddress returns API logger URL
func serverAddress() string {
	envServer := strings.Trim(strings.TrimSpace(os.Getenv("MPS_CLI_SERVER")), "/")
	if len(envServer) == 0 {
		return "https://sca.murphysec.com"
	}
	return envServer
}

var defaultToken string

func SetDefaultToken(token string) {
	defaultToken = token
}
