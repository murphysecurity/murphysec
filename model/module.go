package model

import "github.com/google/uuid"

type PackageManagerType string

const (
	PMMaven     PackageManagerType = "maven"
	PMGoMod     PackageManagerType = "gomod"
	PMNpm       PackageManagerType = "npm"
	PMGradle    PackageManagerType = "gradle"
	PMYarn      PackageManagerType = "yarn"
	PMPython    PackageManagerType = "python"
	PMPip       PackageManagerType = "pip"
	PMComposer  PackageManagerType = "composer"
	PMBundler   PackageManagerType = "bundler"
	PMCocoaPods PackageManagerType = "cocoapods"
	PMPoetry    PackageManagerType = "poetry"
	PmNuget     PackageManagerType = "nuget"
)

type Language string

const (
	Cxx        Language = "C/C++"
	Ruby       Language = "Ruby"
	ObjectiveC Language = "Objective-C"
	PHP        Language = "PHP"
	Go         Language = "Go"
	Java       Language = "Java"
	JavaScript Language = "JavaScript"
	Python     Language = "Python"
	DotNet     Language = "DotNet"
)

type Dependency struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

type Module struct {
	PackageManager PackageManagerType `json:"package_manager"`
	Language       Language           `json:"language"`
	PackageFile    string             `json:"package_file"`
	Name           string             `json:"name"`
	Version        string             `json:"version"`
	FilePath       string             `json:"relative_path"`
	Dependencies   []Dependency       `json:"dependencies,omitempty"`
	RuntimeInfo    interface{}        `json:"runtime_info,omitempty"`
	UUID           uuid.UUID          `json:"uuid"`
}
