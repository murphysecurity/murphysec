package utils

import "os"

var DesensitizedArgs []string

func init() {
	DesensitizedArgs = make([]string, 0, len(os.Args))
	var drop = false
	for _, it := range os.Args {
		if drop {
			drop = false
			DesensitizedArgs = append(DesensitizedArgs, "***")
			continue
		}
		switch it {
		case "--token":
			drop = true
		}
		DesensitizedArgs = append(DesensitizedArgs, it)
	}
}
