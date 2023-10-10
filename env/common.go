package env

import (
	"os"
	"strconv"
	"time"
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
var ScannerScan = false
var ScannerShouldEnableMavenBackupScan = false
var ScannerShouldEnableGradleBackupScan = false
var CommandTimeout time.Duration
var NoWait bool
var envTlsAllowInsecure bool
var CliTlsAllowInsecure bool
var DoNotBuild bool

func init() {
	ctm := os.Getenv("COMMAND_TIMEOUT")
	ct, e := strconv.Atoi(ctm)
	if e != nil || ctm == "" {
		CommandTimeout = time.Second * 25
	} else {
		CommandTimeout = time.Second * time.Duration(ct)
	}
	allowInsecure, e := strconv.ParseBool(os.Getenv("TLS_ALLOW_INSECURE"))
	if allowInsecure && e == nil {
		envTlsAllowInsecure = true
	}
	DoNotBuild, _ = strconv.ParseBool(os.Getenv("DO_NOT_BUILD"))
	DoNotBuild2, _ := strconv.ParseBool(os.Getenv("MPS_DO_NOT_BUILD"))
	DoNotBuild = DoNotBuild || DoNotBuild2
}

func TlsAllowInsecure() bool {
	return CliTlsAllowInsecure || envTlsAllowInsecure
}
