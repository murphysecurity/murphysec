package inlineproperty

import (
	"regexp"
	"strings"
)

type Properties struct {
	m map[string]string
}

func New() *Properties {
	return &Properties{m: map[string]string{}}
}
func (p *Properties) PutIfAbsent(k string, v string) {
	if _, ok := p.m[k]; !ok {
		p.m[k] = v
	}
}
func (p *Properties) Put(k string, v string) {
	p.m[k] = v
}
func (p *Properties) PutMap(m map[string]string) {
	if m == nil {
		return
	}
	for k, v := range m {
		p.Put(k, v)
	}
}

func (p *Properties) Resolve(input string) string {
	return (&ctx{
		m: p.m,
		p: map[string]struct{}{},
	}).resolve(input)
}

type ctx struct {
	m map[string]string
	p map[string]struct{}
}

func (ctx *ctx) resolve(input string) string {
	var raw = pattern.Split(input, -1)
	var refs = make([]string, 0, len(raw))
	var rs = make([]string, len(raw)+len(refs))
	rs = append(rs, raw[0])
	for idx, g := range pattern.FindAllStringSubmatch(input, -1) {
		_, visited := ctx.p[g[1]]
		ctx.p[g[1]] = struct{}{}
		s, ok := ctx.m[g[1]]
		if !ok || visited {
			s = g[0]
		}
		rs = append(rs, s, raw[idx+1])
	}
	s := strings.Join(rs, "")
	if s == input {
		return input
	}
	return ctx.resolve(s)
}

var pattern = regexp.MustCompile(`\$\{([^{}]+)\}`)
