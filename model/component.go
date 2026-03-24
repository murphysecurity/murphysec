package model

// Component is intentionally kept comparable.
//
// Multiple scan and presentation paths use Component as a map key or pass it
// through generic helpers constrained by comparable, for example:
// - module.Module.ComponentList deduplicates via map[Component]struct{}
// - inspector text/code-fragment indexing uses map[Component]...
// - IDEA output joins metadata by map[Component]...
//
// That means fields added here must also remain comparable. In particular, do
// not place slices, maps, or funcs directly on Component unless the related
// call sites are refactored away from map-key equality semantics.
type Component struct {
	CompName    string `json:"comp_name"`
	CompVersion string `json:"comp_version"`
	// SkillFiles is stored behind a pointer so Component stays comparable.
	// The pointed value may contain slices, but the pointer itself is still a
	// valid comparable field for the struct.
	SkillFiles *SkillFiles `json:"skill_files,omitempty"`
	EcoRepo
}

type EcoRepo struct {
	Ecosystem  string `json:"ecosystem"`
	Repository string `json:"repository"`
}

type SkillFile struct {
	RelativePath string       `json:"relative_path"`
	SHA256Hashes []SHA256Hash `json:"sha256_hashes"`
}

// SkillFiles is kept outside Component's value fields because the slice itself
// is not comparable.
type SkillFiles []SkillFile
