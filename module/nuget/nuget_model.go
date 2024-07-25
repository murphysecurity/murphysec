package nuget

import (
	"encoding/xml"

	"github.com/murphysecurity/murphysec/model"
)

var EcoRepo = model.EcoRepo{
	Ecosystem:  "nuget",
	Repository: "",
}

type Project struct {
	XMLName     xml.Name `xml:"Project"`
	PackageRefs []struct {
		Include string `xml:"Include,attr"`
		Version string `xml:"Version,attr"`
	} `xml:"ItemGroup>PackageReference"`
}
type PkgConfig struct {
	XMLName xml.Name `xml:"packages"`
	Package []struct {
		Id                    string `xml:"id,attr"`
		Version               string `xml:"version,attr"`
		DevelopmentDependency bool   `xml:"developmentDependency,attr"`
	} `xml:"package"`
}

type ProjectPackages struct {
	Version    int    `json:"version"`
	Parameters string `json:"parameters"`
	Projects   []struct {
		Path       string `json:"path"`
		Frameworks []struct {
			Framework        string `json:"framework"`
			TopLevelPackages []struct {
				Id               string `json:"id"`
				RequestedVersion string `json:"requestedVersion"`
				ResolvedVersion  string `json:"resolvedVersion"`
			} `json:"topLevelPackages"`
			TransitivePackages []struct {
				Id              string `json:"id"`
				ResolvedVersion string `json:"resolvedVersion"`
			} `json:"transitivePackages"`
		} `json:"frameworks"`
	} `json:"projects"`
}
