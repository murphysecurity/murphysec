package go_mod

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHigherVersion(t *testing.T) {
	tests := []struct {
		name     string
		version1 string
		version2 string
		want     string
	}{
		{
			name:     "second version is higher",
			version1: "v1.2.3",
			version2: "v1.3.0",
			want:     "v1.3.0",
		},
		{
			name:     "first version is higher",
			version1: "v1.3.0",
			version2: "v1.2.3",
			want:     "v1.3.0",
		},
		{
			name:     "semantically equal versions keep first",
			version1: "v1.2.3+incompatible",
			version2: "v1.2.3",
			want:     "v1.2.3+incompatible",
		},
		{
			name:     "release is kept when prerelease is second",
			version1: "v2.0.0",
			version2: "v2.0.0-rc.1",
			want:     "v2.0.0",
		},
		{
			name:     "release is selected when prerelease is first",
			version1: "v2.0.0-rc.1",
			version2: "v2.0.0",
			want:     "v2.0.0",
		},
		{
			name:     "newer pseudo-version is selected when second",
			version1: "v0.0.0-20240101000000-aaaaaaaaaaaa",
			version2: "v0.0.0-20240201000000-bbbbbbbbbbbb",
			want:     "v0.0.0-20240201000000-bbbbbbbbbbbb",
		},
		{
			name:     "newer pseudo-version is kept when first",
			version1: "v0.0.0-20240201000000-bbbbbbbbbbbb",
			version2: "v0.0.0-20240101000000-aaaaaaaaaaaa",
			want:     "v0.0.0-20240201000000-bbbbbbbbbbbb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := higherVersion(tt.version1, tt.version2)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHigherVersionRejectsInvalidVersions(t *testing.T) {
	for _, versions := range [][2]string{
		{"not-a-version", "v1.0.0"},
		{"v1.0.0", "not-a-version"},
	} {
		_, err := higherVersion(versions[0], versions[1])
		assert.Error(t, err)
	}
}
