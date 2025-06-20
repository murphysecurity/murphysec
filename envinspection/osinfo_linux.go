package envinspection

import "github.com/zcalusic/sysinfo"

func getOsInfo() string {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	var vendor = si.OS.Vendor
	var version = si.OS.Version
	switch vendor {
	case "opensuse-leap":
		vendor = "opensuse:leap"
	case "opensuse-leap-micro":
		vendor = "opensuse:leap_micro"
	case "opensuse-tumbleweed":
		vendor = "opensuse:tumbleweed"
		version = ""
	}
	var s = vendor
	if version != "" {
		s += ":" + version
	}
	return s
}
