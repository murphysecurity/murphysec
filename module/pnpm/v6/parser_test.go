package v6

import (
	"context"
	"embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getPureVersion(t *testing.T) {
	assert.Equal(t, "1.0.0", getPureVersion("1.0.0"))
	assert.Equal(t, "13.4.12", getPureVersion("13.4.12(react-dom@18.2.0)(react@18.2.0)"))
	assert.Equal(t, "13.4.12", getPureVersion("/@next/env@13.4.12"))
}

func Test_searchPathIter(t *testing.T) {
	var chunks = []string{
		"/a/b/c(1)(2)n@v",
		"/a/b/c(1)(2)/n@v",
		"/a/b/c/n@v",
		"/a/b/n@v",
		"/a/n@v",
		"/n@v",
		"n@v",
	}
	for s := range searchPathSeq("/a/b/c(1)(2)", "n", "v") {
		assert.Equal(t, chunks[0], s)
		chunks = chunks[1:]
	}
}

//go:embed testdata/*
var testdata embed.FS

func Test_Process(t *testing.T) {
	var entries, e = testdata.ReadDir("testdata")
	assert.NoError(t, e)
	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			var data, e = testdata.ReadFile("testdata/" + entry.Name())
			assert.NoError(t, e)
			tree, e := Process(context.TODO(), data, false)
			assert.NoError(t, e)
			t.Log(tree)
		})
	}

}
