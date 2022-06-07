//go:build windows

package version

import (
	"github.com/murphysecurity/murphysec/logger"
	"golang.org/x/sys/windows/registry"
)

func getOSVersion() string {
	k, e := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if e != nil {
		logger.Err.Println("Open SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion failed.", e.Error())
		return ""
	}
	s, _, e := k.GetStringValue("ProductName")
	if e != nil {
		logger.Err.Println("Read ProductName failed.", e.Error())
	}
	return s
}
