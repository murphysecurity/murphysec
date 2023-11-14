package python

import (
	"github.com/dlclark/regexp2"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/repeale/fp-go"
	"regexp"
	"strings"
)

const pyName = `(?:[A-Za-z0-9][A-Za-z0-9._-]*[A-Za-z0-9]|[A-Za-z0-9])`
const pyVersion = `(?:[A-Za-z0-9_.!-]+)`
const pyVersionOp = `(?<![=!<>])(?:=|<=|==|>=|===)`
const pyVersionSeg = pyVersionOp + `\s*['""]?(` + pyVersion + `)`

var pyVersionSegPattern = regexp2.MustCompile(pyVersionSeg, regexp2.Compiled)
var pyNamePrefixPattern = regexp.MustCompile("^" + pyName)

func parseRequirements(data string) map[string]string {
	var lines []string
	var lineContinuation = false
	for _, s := range fp.Map(func(s string) string { return strings.TrimRight(s, "\r") })(strings.Split(data, "\n")) {
		s = strings.TrimRight(s, "\r")
		var currentLine = strings.TrimSuffix(s, "\\")
		if lineContinuation {
			lines[len(lines)-1] = lines[len(lines)-1] + currentLine
		} else {
			lines = append(lines, currentLine)
		}
		lineContinuation = strings.HasSuffix(s, "\\")
	}
	lines = fp.Map(func(t string) string {
		var i = strings.IndexRune(t, '#')
		if i > -1 {
			return t[:i]
		}
		return t
	})(lines)
	lines = fp.Map(func(t string) string { return strings.TrimSpace(t) })(lines)
	lines = fp.Filter(func(t string) bool { return t != "" })(lines)

	var m = make(map[string]string)
	for _, line := range lines {
		i := strings.IndexRune(line, ';')
		if i > -1 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		name := pyNamePrefixPattern.FindString(line)
		if name == "" {
			continue
		}
		var version string
		var versions []string
		match := must.A(pyVersionSegPattern.FindStringMatch(line))
		for match != nil {
			versions = append(versions, match.GroupByNumber(1).String())
			match = must.A(pyVersionSegPattern.FindNextMatch(match))
		}
		if len(versions) == 1 {
			version = versions[0]
		}
		m[name] = version
	}
	return m
}
