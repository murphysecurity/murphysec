package v5

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getNameVersionFromPath0(t *testing.T) {
	var g = [][3]string{
		{"@babel/plugin-proposal-optional-catch-binding", "7.18.6", "/@babel/plugin-proposal-optional-catch-binding/7.18.6_@babel+core@7.18.6"},
		{"postcss-selector-parser", "6.0.10", "/postcss-selector-parser/6.0.10"},
		{"@typescript-eslint/eslint-plugin", "4.29.0", "/@typescript-eslint/eslint-plugin/4.29.0_48ea228fa0647506aa803d17f48b59f7"},
		{"tsutils", "3.21.0", "/tsutils/3.21.0_typescript@4.3.5"},
	}

	for i, s := range g {
		name, version := getNameVersionFromPath0(s[2])
		assert.Equal(t, s[0], name, "name:", i)
		assert.Equal(t, s[1], version, "version:", i)
	}
}
