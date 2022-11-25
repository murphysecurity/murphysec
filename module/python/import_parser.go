package python

import (
	"regexp"
	"strings"
)

var pyImportPattern1 = regexp.MustCompile(`import\s+(?:[A-Za-z_-][\w.-]*)(?:\s*,\s*(?:[A-Za-z_-][\w.-]*))`)
var pyImportPattern2 = regexp.MustCompile(`from\s+([A-Za-z_-][\w-]*)`)

// returns package name list
func parsePyImport(input string) []string {
	var rs []string
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "import ") {
		// import aa, bb.cc
		for _, it := range strings.Split(strings.TrimPrefix(pyImportPattern1.FindString(input), "import"), ",") {
			it = strings.TrimSpace(it)
			s := strings.Split(it, ".")[0]
			if s != "" {
				rs = append(rs, s)
			}
		}
	}
	if strings.HasPrefix(input, "from ") {
		if m := pyImportPattern2.FindStringSubmatch(input); m != nil {
			rs = append(rs, m[1])
		}
	}
	return rs
}
