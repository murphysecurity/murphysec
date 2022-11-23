package gitinfo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGetSummary(t *testing.T) {
	s, e := filepath.Abs(".")
	assert.NoError(t, e)
	summary, e := _getSummary(t, s)
	assert.NoError(t, e)
	assert.NotEmpty(t, summary.CommitHash)
}

func _getSummary(t *testing.T, dir string) (*Summary, error) {
	summary, e := GetSummary(context.TODO(), dir)
	t.Log(dir, summary, e)
	if e != nil {
		pa := filepath.Dir(dir)
		if pa == dir {
			return nil, e
		}
		return _getSummary(t, pa)
	}
	return summary, e
}
