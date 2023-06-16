package model

import "github.com/murphysecurity/fix-tools/fix"

type ComponentCodeFragment struct {
	Component          Component    `json:"component"`
	CodeFragmentResult fix.Response `json:"code_fragment_result"`
}

type CodeFragment = fix.Preview
