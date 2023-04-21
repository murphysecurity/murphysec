package model

type ComponentCodeFragment struct {
	Component     Component      `json:"component"`
	CodeFragments []CodeFragment `json:"code_fragments,omitempty"`
}

type CodeFragment struct {
	Text         string `json:"text"`
	LineBegin    int    `json:"line_begin"`
	RelativePath string `json:"relative_path"`
}
