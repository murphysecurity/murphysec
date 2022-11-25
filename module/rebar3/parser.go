package rebar3

import (
	"regexp"
	"strings"
)

func parseRebar3TreeOutput(input string) []depNode {
	var t = tokenize(input)
	var roots []depNode
	for {
		if t.get() == nil {
			break
		}
		roots = append(roots, _parse(t, -1)...)
	}
	return roots
}

func _parse(t *tokenizer, indent int) []depNode {
	var rs []depNode
	var currIndent = -1
	for {
		g := t.get()
		if g == nil || g.indent <= indent {
			break
		}
		if currIndent > -1 && g.indent > currIndent {
			rs[len(rs)-1].Children = append(rs[len(rs)-1].Children, _parse(t, currIndent)...)
			continue
		}
		currIndent = g.indent
		rs = append(rs, depNode{
			Name:    g.name,
			Version: g.version,
		})
		t.consume()
	}
	return rs
}

type depNode struct {
	Name     string
	Version  string
	Children []depNode
}

func tokenize(input string) *tokenizer {
	var lp = regexp.MustCompile(`^([ │├─└]*)(.+?)─(.+?)\s`)
	t := &tokenizer{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimRight(line, "\r")
		g := lp.FindStringSubmatch(line)
		if len(g) == 0 {
			continue
		}
		t.g = append(t.g, tokenLine{
			indent:  len(g[1]),
			name:    g[2],
			version: g[3],
		})
	}
	return t
}

type tokenizer struct {
	g []tokenLine
	p int
}

func (t *tokenizer) get() *tokenLine {
	if t.p >= len(t.g) {
		return nil
	}
	return &t.g[t.p]
}

func (t *tokenizer) consume() {
	t.p++
}

type tokenLine struct {
	indent  int
	name    string
	version string
}
