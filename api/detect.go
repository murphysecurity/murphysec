package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
	"net/url"
	"time"
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
	FileHashList   []FileHash     `json:"file_hash_list,omitempty"`
	Language       string         `json:"language"`
	Name           string         `json:"name"`
	PackageFile    string         `json:"package_file"`
	PackageManager string         `json:"package_manager"`
	RelativePath   string         `json:"relative_path"`
	RuntimeInfo    interface{}    `json:"runtime_info,omitempty"`
	Version        string         `json:"version,omitempty"`
	ModuleUUID     uuid.UUID      `json:"module_uuid"`
}

type UserCliDetectInput struct {
	ApiToken           string            `json:"api_token"`
	CliVersion         string            `json:"cli_version"`
	CmdLine            string            `json:"cmd_line"`
	Engine             string            `json:"engine"`
	GitInfo            *VoGitInfo        `json:"git_info"`
	Modules            []VoModule        `json:"modules,omitempty"`
	ProjectName        string            `json:"project_name"`
	TargetAbsPath      string            `json:"target_abs_path"`
	TaskConsumeTime    int               `json:"task_consume_time"`
	TaskInfo           string            `json:"task_info"`
	TaskStartTimestamp int64             `json:"task_start_timestamp"`
	TaskSource         InspectTaskSource `json:"task_type"`
	UserAgent          string            `json:"user_agent"`
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
	LicenseLevelHigh   LicenseLevel = "High"
	LicenseLevelMedium LicenseLevel = "Medium"
	LicenseLevelLow    LicenseLevel = "Low"
)

type VoDetectResponse struct {
	DependenciesCount int `json:"dependencies_count"`
	IssuesCompsCount  int `json:"issues_comps_count"`
	Modules           []struct {
		ModuleId       int       `json:"module_id"`
		ModuleUUID     uuid.UUID `json:"module_uuid"`
		Language       string    `json:"language"`
		PackageManager string    `json:"package_manager"`
		Comps          []struct {
			Solutions []struct {
				Type          string `json:"type"`
				Description   string `json:"description"`
				Compatibility *int   `json:"compatibility"`
			} `json:"solutions"`
			License *struct {
				Spdx  string       `json:"spdx"`
				Level LicenseLevel `json:"level"`
			} `json:"license"`
			IsDirectDependency bool         `json:"is_direct_dependency"`
			CompId             int          `json:"comp_id"`
			CompName           string       `json:"comp_name"`
			CompVersion        string       `json:"comp_version"`
			MinFixedVersion    string       `json:"min_fixed_version"`
			Vuls               []VoVulnInfo `json:"vuls"`
		} `json:"comps"`
	} `json:"modules"`
	DetectStartTimestamp time.Time `json:"detect_start_timestamp"`
	DetectStatus         string    `json:"detect_status"`
	TaskId               string    `json:"task_id"`
}

func SendDetect(input *UserCliDetectInput) (*VoDetectResponse, error) {
	uri := serverAddress() + "/message/v1/access/detect/user_cli"
	logger.Info.Println("Call API:", uri)
	requestData := must.Byte(json.Marshal(input))
	body := bytes.NewReader(requestData)
	request := must.Req(http.NewRequest(http.MethodPost, uri, body))
	request.Header.Set("Content-Type", "application/json")
	// 发送请求
	r, e := client.Do(request)
	if e != nil {
		logger.Err.Println("API request failed.", e.Error())
		if e, ok := e.(*url.Error); ok {
			if e.Timeout() {
				return nil, ErrApiTimeout
			}
		}
		return nil, e
	}
	logger.Info.Println("HTTP request done. Status:", r.StatusCode)
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()
	b, e := io.ReadAll(r.Body)
	if e != nil {
		logger.Err.Println("Read body failed.", e.Error())
		return nil, e
	}
	if r.StatusCode == http.StatusUnauthorized {
		logger.Err.Println("API status:", r.StatusCode, r.Status)
		logger.Err.Println("Invalid token")
		return nil, ErrTokenInvalid
	}
	logger.Debug.Println("API response body")
	logger.Debug.Println(string(b))
	if r.StatusCode != http.StatusOK {
		logger.Err.Println("API status:", r.StatusCode)
		return nil, errors.New(fmt.Sprintf("API status: %d", r.StatusCode))
	}
	v := struct {
		Data VoDetectResponse `json:"data"`
	}{}
	if e := json.Unmarshal(b, &v); e != nil {
		logger.Err.Println("API response body decode failed.", e.Error())
		return nil, e
	}
	return &v.Data, nil
}

func SendDetectHash(input *UserCliDetectInput) (*VoDetectResponse, error) {
	uri := serverAddress() + "/message/v1/access/detect/user_cli_hash"
	logger.Info.Println("Call API:", uri)
	requestData := must.Byte(json.Marshal(input))
	body := bytes.NewReader(requestData)
	request := must.Req(http.NewRequest(http.MethodPost, uri, body))
	request.Header.Set("Content-Type", "application/json")
	// 发送请求
	r, e := client.Do(request)
	if e != nil {
		logger.Err.Println("API request failed.", e.Error())
		return nil, e
	}
	logger.Info.Println("HTTP request done. Status:", r.StatusCode)
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()
	b, e := io.ReadAll(r.Body)
	if e != nil {
		logger.Err.Println("Read body failed.", e.Error())
		return nil, e
	}
	if r.StatusCode == http.StatusUnauthorized {
		logger.Err.Println("API status:", r.StatusCode, r.Status)
		logger.Err.Println("Invalid token")
		return nil, ErrTokenInvalid
	}
	logger.Debug.Println("API response body")
	logger.Debug.Println(string(b))
	if r.StatusCode != http.StatusOK {
		logger.Err.Println("API status:", r.StatusCode)
		return nil, errors.New(fmt.Sprintf("API status: %d", r.StatusCode))
	}
	v := struct {
		Data VoDetectResponse `json:"data"`
	}{}
	if e := json.Unmarshal(b, &v); e != nil {
		logger.Err.Println("API response body decode failed.", e.Error())
		return nil, e
	}
	return &v.Data, nil
}
