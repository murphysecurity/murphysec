package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
	"time"
)

func QueryResult(taskId string) (*TaskScanResponse, error) {
	must.True(taskId != "")
	u := serverAddress() + fmt.Sprintf("/message/v2/access/detect/task_scan?scan_id=%s", taskId)
	logger.Info.Println("Query result at:", u)
	for {
		resp, e := http.Get(u)
		if e != nil {
			return nil, ErrSendRequest
		}
		data, e := readHttpBody(resp)
		if e != nil {
			return nil, e
		}
		if resp.StatusCode == http.StatusOK {
			var r = struct {
				Data TaskScanResponse `json:"data"`
			}{}
			if e := json.Unmarshal(data, &r); e != nil {
				return nil, e
			}
			if !r.Data.Complete {
				logger.Debug.Println("not complete, retry")
				time.Sleep(time.Second * 2)
				continue
			}
			if r.Data.Status == "err" {
				logger.Err.Println("Detect error")
				return nil, errors.New("Detect error")
			}
			if r.Data.Status == "success" {
				logger.Info.Println("Detect Succeeded")
				return &r.Data, nil
			}
			if r.Data.Status == "run" {
				continue
			}
			logger.Warn.Println("Server-side unknown status:", r.Data.Status)
			logger.Debug.Println(string(data))
		}
		return nil, readCommonErr(data, resp.StatusCode)
	}
}

type TaskScanResponse struct {
	Complete          bool `json:"complete"`
	DependenciesCount int  `json:"dependencies_count"`
	IssuesCompsCount  int  `json:"issues_comps_count"`
	Modules           []struct {
		ModuleId       int       `json:"module_id"`
		ModuleUUID     uuid.UUID `json:"module_uuid"`
		Language       string    `json:"language"`
		PackageManager string    `json:"package_manager"`
		Comps          []struct {
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
			Vuls []VoVulnInfo `json:"vuls"`
		} `json:"comps"`
	} `json:"modules"`
	DetectStartTimestamp time.Time `json:"detect_start_timestamp"`
	DetectStatus         string    `json:"detect_status"`
	TaskId               string    `json:"task_id"`
	Status               string    `json:"status"`
}
