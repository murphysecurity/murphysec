package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/util/must"
	"murphysec-cli-simple/util/output"
	"murphysec-cli-simple/util/simplejson"
	"net/http"
	"time"
)

type ScanRequestBody struct {
	CliVersion         string      `json:"cli_version"`
	TaskStatus         int         `json:"task_status"`
	TaskFailureReason  string      `json:"task_failure_reason"`
	TaskType           string      `json:"task_type"`
	OsType             string      `json:"os_type"`
	CmdLine            string      `json:"cmd_line"`
	Plugin             string      `json:"plugin"`
	TaskConsumeTime    int         `json:"task_consume_time"`
	ApiToken           string      `json:"api_token"`
	TaskStartTimestamp int         `json:"task_start_timestamp"`
	ProjectType        string      `json:"project_type"`
	ProjectName        string      `json:"project_name"`
	GitRemoteUrl       string      `json:"git_remote_url"`
	GitBranch          string      `json:"git_branch"`
	TargetPath         string      `json:"target_path"`
	TargetAbsPath      string      `json:"target_abs_path"`
	PackageManager     string      `json:"package_manager"`
	PackageFile        string      `json:"package_file"`
	PackageFilePath    string      `json:"package_file_path"`
	Language           string      `json:"language"`
	TaskResult         interface{} `json:"task_result"`
}

type ScanRequestResponse struct {
}

type ScanResult struct {
	DependenciesCount    int                       `json:"dependencies_count"`
	IssuesCount          int                       `json:"issues_count"`
	DetectStartTimestamp string                    `json:"detect_start_timestamp"`
	DetectConsumeTime    int                       `json:"detect_consume_time"`
	DetectStatus         string                    `json:"detect_status"`
	IssuesLevelCount     ScanResultIssueLevelCount `json:"issues_level_count"`
	TaskId               string                    `json:"task_id"`
	DetectResult         struct {
		VulnInfo []ScanResultVulnInfo `json:"vuln_info"`
	} `json:"detect_result"`
}

type ScanResultVulnInfo struct {
	VulnNo      string                `json:"vuln_no"`
	VulnTitle   string                `json:"vuln_title"`
	Impact      string                `json:"impact"`
	PublishTime string                `json:"publish_time"`
	Influence   int                   `json:"influence"`
	CveId       string                `json:"cve_id"`
	Poc         string                `json:"poc"`
	Cvss        float64               `json:"cvss"`
	Description string                `json:"description"`
	Solution    string                `json:"solution"`
	Source      string                `json:"source"`
	Effect      []ScanResultEffect    `json:"effect"`
	Suggest     string                `json:"suggest"`
	CompName    string                `json:"compName"`
	VulnPath    []string              `json:"vuln_path"`
	References  []ScanResultReference `json:"references"`
}
type ScanResultReference struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ScanResultEffect struct {
	CompNo                string `json:"comp_no"`
	Vendor                string `json:"vendor"`
	Product               string `json:"product"`
	Name                  string `json:"name"`
	AffectVersion         string `json:"affect_version"`
	VersionStartExcluding string `json:"version_start_excluding"`
	VersionStartIncluding string `json:"version_start_including"`
	VersionEndExcluding   string `json:"version_end_excluding"`
	VersionEndIncluding   string `json:"version_end_including"`
	MinFixedVersion       string `json:"min_fixed_version"`
	Language              string `json:"language"`
}

type ScanResultIssueLevelCount struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

func Report(body *ScanRequestBody) (*ScanResult, error) {
	if defaultToken == "" {
		return nil, errors.New("API token not set")
	}
	url := serverAddress() + "/v1/cli/report"
	client := http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(must.Byte(json.Marshal(body))))
	must.Must(err)
	request.Header.Add("Authorization", fmt.Sprintf("token %s", defaultToken))
	output.Debug(fmt.Sprintf("Request: %s", request.RequestURI))
	do, err := client.Do(request)
	output.Debug(fmt.Sprintf("Response: [%d]%s", do.StatusCode, do.Status))
	if err != nil {
		output.Debug(fmt.Sprintf("err: %v", err.Error()))
	}
	//goland:noinspection GoUnhandledErrorResult
	defer do.Body.Close()
	if err != nil {
		return nil, err
	}
	j, err := simplejson.NewFromReader(do.Body)
	if do.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed, statusCode: %d", do.StatusCode)
	}
	if ec := j.Get("code").Int(); ec != 0 {
		return nil, fmt.Errorf("API request failed: %d - %s", ec, j.Get("info").String())
	}
	var r ScanResult
	if e := json.Unmarshal(must.Byte(json.Marshal(j.Get("data"))), &r); e != nil {
		return nil, errors.Wrap(e, "API result unmarshal failed")
	}
	return &r, nil
}
