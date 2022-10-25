package gradle

import (
	"regexp"
	"strings"
)

// catch group: groupId, artifactId, version/variable
var ktsImplPattern = regexp.MustCompile(`(?:implementation|runtimeOnly)\("([\w.-]+):([\w.-]+):(\$\w+|[\w.-]+)"\)`)

// catch group: groupId, artifactId, version
var groovyImplPattern = regexp.MustCompile(`(?:implementation|runtimeOnly|compile)['"]([\w.-]+):([\w.-]+):([\w.-]+)['"]`)

// catch group: variableIdentifier, string literal value
var ktsVariablePattern = regexp.MustCompile(`va[lr]\s+(\w)\s*=\s"(.+?)"`)

var replacePattern = regexp.MustCompile(`[\r\s]+`)
var commentPattern = regexp.MustCompile(`//.+$`)

func parseGradleGroovy(input string) []DepElement {
	var rs []DepElement
	for _, it := range strings.Split(input, "\n") {
		// remove all space character, line-break and single line comment
		it = commentPattern.ReplaceAllString(replacePattern.ReplaceAllString(it, ""), "")
		m := groovyImplPattern.FindStringSubmatch(it)
		if m == nil {
			continue
		}
		rs = append(rs, DepElement{
			GroupId:    m[1],
			ArtifactId: m[2],
			Version:    m[3],
		})
	}
	return rs
}

func parseGradleKts(input string) []DepElement {
	var rs []DepElement
	variableMap := map[string]string{}
	for _, s := range strings.Split(input, "\n") {
		if m := ktsVariablePattern.FindStringSubmatch(s); m != nil {
			variableMap[m[1]] = m[2]
			continue
		}
		s = commentPattern.ReplaceAllString(replacePattern.ReplaceAllLiteralString(s, ""), "")
		m := ktsImplPattern.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		v := m[3]
		if strings.HasPrefix(v, "$") {
			if variableMap[v] == "" {
				continue
			}
			v = variableMap[v]
		}
		rs = append(rs, DepElement{
			GroupId:    m[1],
			ArtifactId: m[2],
			Version:    v,
		})
	}
	return rs
}
