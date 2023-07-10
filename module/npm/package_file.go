package npm

import (
	"encoding/json"
	"fmt"
)

type pkgFile struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
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
