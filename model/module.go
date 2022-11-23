package model

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
