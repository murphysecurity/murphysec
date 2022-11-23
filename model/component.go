package model

type Component struct {
	CompName    string `json:"comp_name"`
	CompVersion string `json:"comp_version"`
	EcoRepo
}

type EcoRepo struct {
	Ecosystem  string `json:"ecosystem"`
	Repository string `json:"repository"`
}
