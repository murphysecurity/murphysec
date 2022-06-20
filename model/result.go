package model

type TaskScanResponse struct {
	Complete          bool `json:"complete"`
	DependenciesCount int  `json:"dependencies_count"`
	IssuesCompsCount  int  `json:"issues_comps_count"`
	ProjectScore      int  `json:"project_score"`
	SurpassScore      int  `json:"surpass_score"`
	Modules           []struct {
		ModuleId       int    `json:"module_id"`
		Language       string `json:"language"`
		PackageManager string `json:"package_manager"`
		Comps          []struct {
			MinFixedInfo []struct {
				Name               string `json:"name"`
				OldVersion         string `json:"old_version"`
				NewVersion         string `json:"new_version"`
				SecurityScore      int    `json:"security_score"`
				CompatibilityScore int    `json:"compatibility_score"`
			} `json:"min_fixed_info,omitempty"`
			IsDirectDependency bool   `json:"is_direct_dependency"`
			CompId             int    `json:"comp_id"`
			CompName           string `json:"comp_name"`
			CompVersion        string `json:"comp_version"`
			MinFixedVersion    string `json:"min_fixed_version"`
			License            *struct {
				Level LicenseLevel `json:"level"`
				Spdx  string       `json:"spdx"`
			} `json:"license,omitempty"`
			Solutions []struct {
				Compatibility *int   `json:"compatibility,omitempty"`
				Description   string `json:"description"`
				Type          string `json:"type,omitempty"`
			} `json:"solutions,omitempty"`
			Vuls         []VoVulnInfo `json:"vuls"`
			FixType      string       `json:"fix_type"`
			CompSecScore int          `json:"comp_sec_score"`
		} `json:"comps"`
	} `json:"modules"`
	TaskId           string `json:"task_id"`
	Status           string `json:"status"`
	InspectReportUrl string `json:"inspect_report_url"`
}

type VoVulnInfo struct {
	CveId           string        `json:"cve_id"`
	Description     string        `json:"description"`
	Level           VulnLevelType `json:"level"`
	Influence       int           `json:"influence"`
	Poc             bool          `json:"poc"`
	PublishTime     int           `json:"publish_time"`
	AffectedVersion string        `json:"affected_version"`
	MinFixedVersion string        `json:"min_fixed_version"`
	References      []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"references"`
	Solutions []struct {
		Type          string `json:"type"`
		Description   string `json:"description"`
		Compatibility int    `json:"compatibility"`
	} `json:"solutions"`
	SuggestLevel SuggestLevel `json:"suggest_level"`
	VulnNo       string       `json:"vuln_no"`
	VulnPath     []string     `json:"vuln_path"`
	Title        string       `json:"title"`
}

type SuggestLevel string

const (
	SuggestLevelOptional        SuggestLevel = "Optional"
	SuggestLevelRecommend       SuggestLevel = "Recommend"
	SuggestLevelStrongRecommend SuggestLevel = "StrongRecommend"
)

type VulnLevelType string

const (
	VulnLevelCritical VulnLevelType = "Critical"
	VulnLevelHigh     VulnLevelType = "High"
	VulnLevelMedium   VulnLevelType = "Medium"
	VulnLevelLow      VulnLevelType = "Low"
)

type LicenseLevel string

const (
	LicenseLevelLow    LicenseLevel = "Low"
	LicenseLevelMedium LicenseLevel = "Medium"
	LicenseLevelHigh   LicenseLevel = "High"
)
