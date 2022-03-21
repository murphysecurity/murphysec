package base

import (
	"fmt"
	"github.com/google/uuid"
	"murphysec-cli-simple/api"
	"regexp"
	"strings"
)

type Module struct {
	PackageManager string       `json:"package_manager"`
	Language       string       `json:"language"`
	PackageFile    string       `json:"package_file"`
	Name           string       `json:"name"`
	Version        string       `json:"version"`
	RelativePath   string       `json:"relative_path"`
	Dependencies   []Dependency `json:"dependencies,omitempty"`
	RuntimeInfo    interface{}  `json:"runtime_info,omitempty"`
	UUID           uuid.UUID    `json:"uuid"`
}

func (m Module) ApiVo() *api.VoModule {
	r := &api.VoModule{
		Dependencies:   mapVoDependency(m.Dependencies),
		Language:       m.Language,
		Name:           m.Name,
		PackageFile:    m.PackageFile,
		PackageManager: m.PackageManager,
		RelativePath:   m.RelativePath,
		RuntimeInfo:    m.RuntimeInfo,
		Version:        m.Version,
		ModuleType:     "version",
		ModuleUUID:     m.UUID,
	}
	return r
}

type Dependency struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

var paddingPattern = regexp.MustCompile("^[\\r\\n\\t ]*|[\\r\\n\\t ]*$")

func trimPadding(s string) string {
	return paddingPattern.ReplaceAllString(s, "")
}
func mapVoDependency(d []Dependency) []api.VoDependency {
	r := make([]api.VoDependency, 0)
	for _, it := range d {
		r = append(r, api.VoDependency{
			Name:         trimPadding(it.Name),
			Version:      trimPadding(it.Version),
			Dependencies: mapVoDependency(it.Dependencies),
		})
	}
	return r
}

type Inspector interface {
	fmt.Stringer
	Version() string
	CheckDir(dir string) bool
	Inspect(dir string) ([]Module, error)
	PackageManagerType() PackageManagerType
}

type PackageManagerType string

const (
	PMMaven  PackageManagerType = "maven"
	PMGoMod  PackageManagerType = "gomod"
	PMNpm    PackageManagerType = "npm"
	PMGradle PackageManagerType = "gradle"
	PMYarn   PackageManagerType = "yarn"
	PMPython PackageManagerType = "python"
)

func PackageManagerTypeOfName(name string) PackageManagerType {
	switch PackageManagerType(strings.ToLower(name)) {
	case PMNpm, PMGoMod, PMMaven, PMGradle, PMYarn, PMPython:
		return PackageManagerType(strings.ToLower(name))
	default:
		panic("wtf?")
	}
}
