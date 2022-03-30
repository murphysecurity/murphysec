package api

import (
	"github.com/google/uuid"
	"murphysec-cli-simple/utils/must"
	"time"
)

type VoGitInfo struct {
	Commit        string    `json:"commit"`
	GitRef        string    `json:"git_ref"`
	GitRemoteUrl  string    `json:"git_remote_url"`
	CommitMessage string    `json:"commit_message"`
	CommitEmail   string    `json:"commit_email"`
	CommitTime    time.Time `json:"commit_time"`
}

type VoDependency struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Dependencies []VoDependency `json:"dependencies"`
}

type VoFileHash struct {
	Hash string `json:"hash"`
}
type VoModule struct {
	Dependencies   []VoDependency `json:"dependencies,omitempty"`
	FileHashList   []VoFileHash   `json:"file_hash_list"`
	Language       string         `json:"language,omitempty"`
	Name           string         `json:"name,omitempty"`
	PackageFile    string         `json:"package_file,omitempty"`
	PackageManager string         `json:"package_manager,omitempty"`
	RelativePath   string         `json:"relative_path,omitempty"`
	RuntimeInfo    interface{}    `json:"runtime_info,omitempty"`
	Version        string         `json:"version,omitempty"`
	ModuleUUID     uuid.UUID      `json:"module_uuid,omitempty"`
	ModuleType     string         `json:"module_type"`
}

type FileHash struct {
	Hash string `json:"hash"`
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

type SendDetectRequest struct {
	TaskInfo string     `json:"task_info"`
	ApiToken string     `json:"api_token"`
	Modules  []VoModule `json:"modules"`
}

func SendDetect(input *SendDetectRequest) error {
	must.True(input != nil)
	input.ApiToken = C.Token
	req := C.PostJson("/message/v2/access/detect/user_cli", input)
	return C.DoJson(req, nil)
}
