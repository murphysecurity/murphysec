package envinspection

import (
	"io"
	"os"
	"strconv"
	"strings"
)

var candidateOsReleasePath = []string{"/etc/os-release", "/usr/lib/os-release"}

func readLinuxKernelVersion() string {
	data, e := os.ReadFile("/proc/version")
	if e != nil {
		return ""
	}
	m := strings.SplitN(string(data), " ", 4)
	if len(m) < 4 {
		return ""
	}
	return m[2]
}

func readOsRelease() map[string]string {
	for _, it := range candidateOsReleasePath {
		if rs, e := parseOsRelease(it); e == nil {
			return rs
		}
	}
	return map[string]string{}
}

func parseOsRelease(p string) (map[string]string, error) {
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	data, e := io.ReadAll(io.LimitReader(f, 1*1024*1024))
	if e != nil {
		return nil, e
	}
	rs := map[string]string{}
	for _, s := range strings.Split(string(data), "\n") {
		s = strings.TrimSpace(s)
		if s == "" || s[0] == '#' {
			continue
		}
		line := strings.SplitN(s, "=", 2)
		if len(line) != 2 {
			continue
		}
		var key, value = line[0], line[1]
		if v, e := strconv.Unquote(key); e == nil {
			key = v
		}
		if value != "" && value[0] == '"' {
			if v, e := strconv.Unquote(value); e == nil {
				value = v
			}
		}
		rs[key] = value

	}
	return rs, nil
}
