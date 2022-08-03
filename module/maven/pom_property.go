package maven

import (
	"regexp"
	"strings"
)

type properties struct {
	m map[string]string
}

func newProperties() *properties {
	return &properties{m: map[string]string{}}
}
func (p *properties) PutIfAbsent(k string, v string) {
	if _, ok := p.m[k]; !ok {
		p.m[k] = v
	}
}
func (p *properties) Put(k string, v string) {
	p.m[k] = v
}
func (p *properties) PutMap(m map[string]string) {
	if m == nil {
		return
	}
	for k, v := range m {
		p.Put(k, v)
	}
}

func (p *properties) Resolve(input string) string {
	return (&propertiesCtx{
		m: p.m,
		p: map[string]struct{}{},
	})._resolve(input)
}

type propertiesCtx struct {
	m map[string]string
	p map[string]struct{}
}

func (ctx *propertiesCtx) _resolve(input string) string {
	var raw = inlineProperties.Split(input, -1)
	var refs = make([]string, 0, len(raw))
	var rs = make([]string, 0, len(raw)+len(refs))
	rs = append(rs, raw[0])
	for idx, g := range inlineProperties.FindAllStringSubmatch(input, -1) {
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
	return ctx._resolve(s)
}

var inlineProperties = regexp.MustCompile(`\$\{([^{}]+)\}`)
