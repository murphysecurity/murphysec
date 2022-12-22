package python

import "regexp"

var __dockerFilePipInstallPattern = regexp.MustCompile(`pip\d?\s+install.*?\s-r\s+([^\s&|;"']+)`)

func parseDockerfilePipInstall(input string) []string {
	var r []string
	for _, match := range __dockerFilePipInstallPattern.FindAllStringSubmatch(input, -1) {
		r = append(r, match[1])
	}
	return r
}
