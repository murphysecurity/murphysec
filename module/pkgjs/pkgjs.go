// Package pkgjs contains common utils to processing package.json
package pkgjs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"unsafe"
)

type Pkg struct {
	Name            string            `json:"name"`
	Version         string            `json:"version,omitempty"`
	License         string            `json:"license,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}

func (p Pkg) String() string {
	data, e := json.MarshalIndent(p, "", "  ")
	if e != nil {
		return "[JSON marshal error]"
	}
	return unsafe.String(&data[0], len(data))
}

const Filename = "package.json"

// ReadFile parse the specified file as package.json
//
// Errors: *os.PathError or wrapped other error
func ReadFile(path string) (*Pkg, error) {
	f, e := os.Open(path)
	if e != nil {
		if os.IsNotExist(e) {
			return nil, e
		}
		return nil, wrap(e)
	}
	defer func() { _ = f.Close() }()
	var dec = json.NewDecoder(f)
	var decoded Pkg
	e = dec.Decode(&decoded)
	if e != nil {
		return nil, wrap(e)
	}
	return &decoded, nil
}

// ReadDir parse the package.json file in the directory
//
// Errors: *os.PathError or wrapped other error
func ReadDir(path string) (*Pkg, error) {
	return ReadFile(filepath.Join(path, Filename))
}
