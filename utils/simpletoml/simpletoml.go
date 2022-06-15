package simpletoml

import "github.com/pelletier/go-toml/v2"

type TOML struct {
	V any
}

func UnmarshalTOML(data []byte) (*TOML, error) {
	t := &TOML{}
	return t, toml.Unmarshal(data, &t.V)
}

func (t TOML) Get(path ...string) *TOML {
	var cur = t.V
	for _, it := range path {
		if m, ok := cur.(map[string]any); ok {
			cur = m[it]
		} else {
			return &TOML{}
		}
	}
	return &TOML{V: cur}
}

func (t TOML) String(defaultValue ...string) string {
	if len(defaultValue) > 1 {
		panic("defaultValue.len > 1")
	}
	if s, ok := t.V.(string); ok {
		return s
	} else {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
}

func (t TOML) TOMLArray() (rs []TOML) {
	if s, ok := t.V.([]any); ok {
		for _, it := range s {
			rs = append(rs, TOML{V: it})
		}
	}
	return
}
