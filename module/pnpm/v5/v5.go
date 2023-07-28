package v5

import "gopkg.in/yaml.v3"

type Pkg struct {
	Name         string            `json:"name" yaml:"name"`
	Version      string            `json:"version" yaml:"version"`
	Dependencies map[string]string `json:"dependencies" yaml:"dependencies"`
	Dev          bool              `json:"dev" yaml:"dev"`
}

type Importer struct {
	Dependencies    map[string]string `json:"dependencies" yaml:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies" yaml:"devDependencies"`
}

type Lockfile struct {
	Importers       map[string]*Importer `json:"importers" yaml:"importers"`
	Dependencies    map[string]string    `json:"dependencies" yaml:"dependencies"`
	DevDependencies map[string]string    `json:"devDependencies" yaml:"devDependencies"`
	Packages        map[string]*Pkg      `json:"packages" yaml:"packages"`
}

func ParseLockfile(data []byte) (*Lockfile, error) {
	var lockfile Lockfile
	if e := yaml.Unmarshal(data, &lockfile); e != nil {
		return nil, e
	}
	return &lockfile, nil
}
