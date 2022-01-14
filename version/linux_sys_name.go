//go:build linux

package version

import (
	"io/ioutil"
	"murphysec-cli-simple/logger"
	"regexp"
)

func getOSVersion() string {
	data, e := ioutil.ReadFile("/etc/os-release")
	if e != nil {
		logger.Err.Println("get /etc/os-release failed.", e.Error())
	} else {
		if m := regexp.MustCompile("PRETTY_NAME\\s*=\\s*(?:\\\"([^\\\"\\n]+)\\\"|([^\\n]+))").FindStringSubmatch(string(data)); m != nil {
			if m[1] == "" {
				return m[2]
			}
			return m[1]
		}
	}
	return ""
}
