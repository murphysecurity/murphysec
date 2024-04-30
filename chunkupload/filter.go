package chunkupload

import (
	"io/fs"
	"strings"
)

type FilterVote int

const (
	_ FilterVote = iota
	FilterAdd
	FilterSkip
	FilterSkipDir
)

type Filter func(path string, entry fs.DirEntry) (FilterVote, error)

var DiscardDot Filter = func(path string, entry fs.DirEntry) (FilterVote, error) {
	var name = entry.Name()
	if strings.HasPrefix(name, ".") || name == "node_modules" {
		if entry.IsDir() {
			return FilterSkipDir, nil
		} else {
			return FilterSkip, nil
		}
	}
	return FilterAdd, nil
}
var uploadAll Filter = func(path string, entry fs.DirEntry) (FilterVote, error) { return FilterAdd, nil }
