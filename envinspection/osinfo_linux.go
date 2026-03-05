package envinspection

import (
	"os"
	"regexp"
	"strings"

	"github.com/zcalusic/sysinfo"
)

func getOsInfo() string {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	var vendor = si.OS.Vendor
	var version = si.OS.Version
	// centos/rhel 鍏煎琛ヤ竵
	if vendor == "" || vendor == "unknown" {
		if content, err := os.ReadFile("/etc/redhat-release"); err == nil {
			line := strings.ToLower(string(content))
			var re = regexp.MustCompile(`([a-zA-Z ]+) release ([0-9.]+)`)
			if m := re.FindStringSubmatch(line); len(m) == 3 {
				v := strings.TrimSpace(m[1])
				switch {
				case strings.Contains(v, "centos"):
					vendor = "centos"
				case strings.Contains(v, "red hat enterprise linux server"):
					vendor = "rhel"
				default:
					vendor = strings.ReplaceAll(v, " ", "_")
				}
				// Keep major OS version only
				ver := m[2]
				if idx := strings.Index(ver, "."); idx > 0 {
					version = ver[:idx]
				} else {
					version = ver
				}
			}
		}
	}
	switch vendor {
	case "opensuse-leap":
		vendor = "opensuse:leap"
	case "opensuse-leap-micro":
		vendor = "opensuse:leap_micro"
	case "opensuse-tumbleweed":
		vendor = "opensuse:tumbleweed"
		version = ""
	case "baidulinux":
		vendor = "baidu"
	}
	var s = vendor
	if version != "" {
		s += ":" + version
	}
	return s
}
