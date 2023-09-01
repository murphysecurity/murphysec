package shared

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNameFromPath(t *testing.T) {
	assert.Equal(t, "mdast-util-mdx-expression", ignoreError(GetNameFromPath("/mdast-util-mdx-expression@1.3.2")))
	assert.Equal(t, "android-arm64", ignoreError(GetNameFromPath("/@esbuild/android-arm64@0.17.19")))
	assert.Equal(t, "react-markdown", ignoreError(GetNameFromPath("/react-markdown@8.0.7(@types/react@18.2.6)(react@18.2.0)")))
}

func TestGetVersionFromPath(t *testing.T) {
	assert.Equal(t, "1.3.2", ignoreError(GetVersionFromPath("/mdast-util-mdx-expression@1.3.2")))
	assert.Equal(t, "0.17.19", ignoreError(GetVersionFromPath("/@esbuild/android-arm64@0.17.19")))
	assert.Equal(t, "8.0.7", ignoreError(GetVersionFromPath("/react-markdown@8.0.7(@types/react@18.2.6)(react@18.2.0)")))
}

func Test_trimParenthesesFromPath(t *testing.T) {
	assert.Equal(t, "/react-markdown@8.0.7", trimParenthesesFromPath("/react-markdown@8.0.7(@types/react@18.2.6)(react@18.2.0)"))
}
