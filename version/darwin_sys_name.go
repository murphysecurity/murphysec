//go:build darwin

package version

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"murphysec-cli-simple/logger"
)

func getOSVersion() string {
	b, e := ioutil.ReadFile("/System/Library/CoreServices/SystemVersion.plist")
	if e != nil {
		logger.Err.Println("Open /System/Library/CoreServices/SystemVersion.plist failed.", e.Error())
		return ""
	}
	var p plist
	if e := xml.Unmarshal(b, &p); e != nil {
		logger.Err.Println("plist decode failed,", e.Error())
		return ""
	}
	var id = 0
	for i, s := range p.Dict.Key {
		if s == "ProductVersion" {
			id = i
		}
	}
	if id < len(p.Dict.String) {
		return fmt.Sprintf("macOS %s", p.Dict.String[id])
	}
	return ""
}

type plist struct {
	XMLName xml.Name `xml:"plist"`
	Dict    struct {
		Key    []string `xml:"key"`
		String []string `xml:"string"`
	} `xml:"dict"`
}
