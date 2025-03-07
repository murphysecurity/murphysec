package envinspection

import "github.com/zcalusic/sysinfo"

func getOsInfo() string {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	return si.OS.Vendor + ":" + si.OS.Version
}
