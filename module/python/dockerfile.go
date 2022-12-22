package python

import (
	"regexp"
	"strings"
)

var __dockerFilePipInstallPattern = regexp.MustCompile(`pip\d?\s+install.*?\s-r\s+([^\s&|;"']+)`)

func parseDockerfilePipInstall(input string) []string {
	var r []string
	for _, match := range __dockerFilePipInstallPattern.FindAllStringSubmatch(input, -1) {
		s := strings.TrimSpace(match[1])
		if s == "" {
			continue
		}
		r = append(r, s)
	}
	return r
}
