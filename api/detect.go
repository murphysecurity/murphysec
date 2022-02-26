package api

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
)

type VoGitInfo struct {
	Commit       string `json:"commit"`
	GitRef       string `json:"git_ref"`
	GitRemoteUrl string `json:"git_remote_url"`
}

type VoDependency struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Dependencies []VoDependency `json:"dependencies"`
}

type VoModule struct {
	Dependencies   []VoDependency `json:"dependencies,omitempty"`
	Hash           []string       `json:"hash,omitempty"`
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
	u := serverAddress() + "/message/v2/access/detect/user_cli"
	input.ApiToken = conf.APIToken()
	body := must.Byte(json.Marshal(input))
	resp, e := http.Post(u, "application/json", bytes.NewReader(body))
	if e != nil {
		logger.Warn.Println("send request failed.", e.Error())
		return ErrSendRequest
	}
	if resp.StatusCode == 200 {
		logger.Info.Println("SendDetect succeeded")
		return nil
	}
	data, e := readHttpBody(resp)
	if e != nil {
		return e
	}
	return readCommonErr(data, resp.StatusCode)
}
