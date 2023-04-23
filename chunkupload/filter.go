package chunkupload

import "io/fs"

type FilterVote int

const (
	_ FilterVote = iota
	FilterAdd
	FilterSkip
	FilterSkipDir
)

type Filter func(path string, entry fs.DirEntry) (FilterVote, error)

var uploadAll Filter = func(path string, entry fs.DirEntry) (FilterVote, error) { return FilterAdd, nil }
