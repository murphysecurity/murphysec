package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// serverAddress returns API logger URL
func serverAddress() string {
	envServer := strings.Trim(strings.TrimSpace(os.Getenv("MPS_CLI_SERVER")), "/")
	if len(envServer) == 0 {
		return "https://sca.murphysec.com"
	}
	return envServer
}

var client *http.Client

func init() {
	c := new(http.Client)
	c.Timeout = time.Second * 30
	i, e := strconv.Atoi(os.Getenv("API_TIMEOUT"))
	if e == nil && i > 0 {
		c.Timeout = time.Duration(int64(time.Second) * int64(i))
	}
	client = c
}
func getClient() *http.Client {
	return client
}
