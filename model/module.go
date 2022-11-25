package model

import "github.com/murphysecurity/murphysec/utils"

type Module struct {
	ModuleName     string           `json:"module_name"`
	ModuleVersion  string           `json:"module_version"`
	ModulePath     string           `json:"module_path"`
	PackageManager string           `json:"package_manager"`
	Dependencies   []DependencyItem `json:"dependencies,omitempty"`
	ScanStrategy   ScanStrategy     `json:"scan_strategy"`
}

func (m Module) String() string {
	var s = "[" + m.PackageManager + "]" + m.ModuleName
	if m.ModuleVersion != "" {
		s += "@" + m.ModuleVersion
	}
	return s
}

func (m Module) IsZero() bool {
	return len(m.Dependencies) == 0 && m.ModuleName == "" && m.ModuleVersion == ""
}

func (m Module) ComponentList() []Component {
	var r = make(map[Component]struct{})
	__componentList(m.Dependencies, r)
	return utils.KeysOfMap(r)
}

func __componentList(deps []DependencyItem, m map[Component]struct{}) {
	for _, dep := range deps {
		m[dep.Component] = struct{}{}
		__componentList(dep.Dependencies, m)
	}
}
