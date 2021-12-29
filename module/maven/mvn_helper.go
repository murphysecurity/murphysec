package maven

import (
	"bufio"
	"io"
	"regexp"
)

var modulePattern = regexp.MustCompile("digraph +\\\"(.+?):(.+?):.+?:(.+?)\\\" .*\\{")
var depPattern = regexp.MustCompile("\\\"(?:(.+?):(.+?):.+?:(.+?))[:\\\"].*?->\\s+\\\"(?:(.+?):(.+?):.+?:(.+?))[\\\":]")

const _MaxLineSize = 32 * 1024

func parseOutput(reader io.Reader) map[Coordination]map[Coordination][]Coordination {
	rs := map[Coordination]map[Coordination][]Coordination{}

	input := bufio.NewScanner(reader)
	input.Split(bufio.ScanLines)
	input.Buffer(make([]byte, _MaxLineSize), _MaxLineSize)

	var currentModule Coordination

	for input.Scan() {
		line := input.Text()
		if m := modulePattern.FindStringSubmatch(line); m != nil {
			// module matched
			currentModule = Coordination{
				GroupId:    m[1],
				ArtifactId: m[2],
				Version:    m[3],
			}
			continue
		}
		if m := depPattern.FindStringSubmatch(line); m != nil {
		}
	}

	return rs
}
