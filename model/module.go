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
	PmConan     PackageManagerType = "conan"
	PmRebar3    PackageManagerType = "rebar3"
	PmCargo     PackageManagerType = "cargo"
	PmIvy       PackageManagerType = "ivy"
	PmSbt       PackageManagerType = "sbt"
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
	Erlang     Language = "Erlang"
	Rust       Language = "Rust"
	Scala      Language = "Scala"
)

type ScanStrategy string

const (
	ScanStrategyNormal ScanStrategy = "Normal"
	ScanStrategyBackup ScanStrategy = "Backup"
)

type Dependency struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

type Module struct {
	PackageManager PackageManagerType `json:"package_manager"`
	Language       Language           `json:"language"`
	Name           string             `json:"name"`
	Version        string             `json:"version"`
	RelativePath   string             `json:"relative_path"`
	Dependencies   []Dependency       `json:"dependencies,omitempty"`
	RuntimeInfo    interface{}        `json:"runtime_info,omitempty"`
	UUID           uuid.UUID          `json:"uuid"`
	ScanStrategy   ScanStrategy       `json:"scan_strategy"`
}
