package npm

import (
	"encoding/json"
	"fmt"
	"sort"
)

type pkgFile struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (p pkgFile) DependenciesEntries() [][2]string {
	var r [][2]string
	for n, v := range p.Dependencies {
		r = append(r, [2]string{n, v})
	}
	sortEntries(r)
	return r
}

func (p pkgFile) DevDependenciesEntries() [][2]string {
	var r [][2]string
	for n, v := range p.DevDependencies {
		r = append(r, [2]string{n, v})
	}
	sortEntries(r)
	return r
}

func sortEntries(input [][2]string) {
	sort.Slice(input, func(i, j int) bool {
		if input[i][0] == input[j][0] {
			return input[i][1] < input[j][1]
		}
		return input[i][0] < input[j][0]
	})
}

func parsePkgFile(data []byte) (*pkgFile, error) {
	var e error
	var r pkgFile
	e = json.Unmarshal(data, &r)
	if e != nil {
		return nil, fmt.Errorf("parsing package file: bad format, %w", e)
	}
	return &r, nil
}
