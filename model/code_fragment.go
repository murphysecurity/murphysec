package model

import "github.com/murphysecurity/fix-tools/fix"

type ComponentCodeFragment struct {
	Component     Component      `json:"component"`
	CodeFragments []CodeFragment `json:"code_fragments,omitempty"`
}

type CodeFragment = fix.Preview
