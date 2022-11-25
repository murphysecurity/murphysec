package gradle

import (
	"fmt"
	"regexp"
	"strings"
)

type blockParser struct {
	lines   []string
	linePtr int
}

func (b *blockParser) eof() bool {
	return b.linePtr >= len(b.lines)
}

func (b *blockParser) consume() {
	b.linePtr++
}

// Gets the value of the symbol at offset i from the current position.
// When i==1, this method returns the value of the current symbol in the stream (which is the next symbol to be consumed).
// When i==-1, this method returns the value of the previously read symbol in the stream.
func (b *blockParser) la(i int) string {
	if i > 0 {
		return b.lines[b.linePtr+i-1]
	} else if i < 0 {
		return b.lines[b.linePtr+i]
	} else {
		panic("invalid")
	}
}

func __calcIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, "|+-\\/ "))
}

func (b *blockParser) _parse() []DepElement {
	var rs []DepElement
	lastIndent := __calcIndent(b.la(1))
	for !b.eof() {
		currIndent := __calcIndent(b.la(1))
		if lastIndent == currIndent {
			item := parseDepElement(b.la(1))
			if item == nil {
				item = &DepElement{ArtifactId: b.la(1)}
			}
			rs = append(rs, *item)
			b.consume()
			continue
		}
		if lastIndent > currIndent {
			return rs
		}
		if lastIndent < currIndent {
			rs[len(rs)-1].Children = b._parse()
		}
	}
	return rs
}

var __parseDepElementPattern1 = regexp.MustCompile(`^([A-Za-z0-9\.-]+)\:([A-Za-z0-9\.-]+)(?:\:([A-Za-z0-9\.-]+))?(?: *-> *([A-Za-z0-9\.-]+))?`)
var __parseDepElementPattern2 = regexp.MustCompile(`^project ([A-Za-z0-9_.:-]+)`)

func parseDepElement(s string) *DepElement {
	s = strings.TrimLeft(s, "+- |\\/")
	if m := __parseDepElementPattern1.FindStringSubmatch(s); m != nil {
		d := &DepElement{
			GroupId:    m[1],
			ArtifactId: m[2],
			Version:    m[3],
		}
		if m[4] != "" {
			d.Version = m[4]
		}
		return d
	}
	if m := __parseDepElementPattern2.FindStringSubmatch(s); m != nil {
		return &DepElement{ArtifactId: m[1]}
	}
	return nil
}

type DepElement struct {
	GroupId    string       `json:"group_id"`
	ArtifactId string       `json:"artifact_id"`
	Version    string       `json:"version"`
	Children   []DepElement `json:"children,omitempty"`
}

func (d DepElement) CompName() string {
	return fmt.Sprintf("%s:%s", d.GroupId, d.ArtifactId)
}
