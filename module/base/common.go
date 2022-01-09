package base

import (
	"fmt"
	"strings"
)

type Module struct {
	PackageManager string       `json:"package_manager"`
	Language       string       `json:"language"`
	PackageFile    string       `json:"package_file"`
	Name           string       `json:"name"`
	Version        string       `json:"version"`
	RelativePath   string       `json:"relative_path"`
	Dependencies   []Dependency `json:"dependencies"`
	RuntimeInfo    interface{}  `json:"runtime_info"`
}

type Dependency struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Dependencies []Dependency `json:"dependencies"`
}

type CommonScanFunc func(dir string) ([]Module, error)
type CommonCheckFunc func(dir string) bool

type Inspector interface {
	fmt.Stringer
	Version() string
	CheckDir(dir string) bool
	Inspect(dir string) ([]Module, error)
	PackageManagerType() PackageManagerType
}

type PackageManagerType string

const (
	PMMaven PackageManagerType = "maven"
	PMGoMod PackageManagerType = "gomod"
	PMNpm   PackageManagerType = "npm"
)

func PackageManagerTypeOfName(name string) PackageManagerType {
	switch PackageManagerType(strings.ToLower(name)) {
	case PMNpm, PMGoMod, PMMaven:
		return PackageManagerType(strings.ToLower(name))
	default:
		panic("wtf?")
	}
}
