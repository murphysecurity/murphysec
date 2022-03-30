package api

import (
	"fmt"
	"murphysec-cli-simple/utils/must"
	"time"
)

func QueryResult(taskId string) (*TaskScanResponse, error) {
	must.True(taskId != "")
	for {
		var r = struct {
			Data TaskScanResponse `json:"data"`
		}{}
		httpReq := C.GET(fmt.Sprintf("/message/v2/access/detect/task_scan?scan_id=%s", taskId))
		if e := C.DoJson(httpReq, &r); e != nil {
			return nil, e
		}
		if !r.Data.Complete {
			time.Sleep(time.Second * 2)
			continue
		}
		return &r.Data, nil
	}
}

type TaskScanResponse struct {
	Complete          bool `json:"complete"`
	DependenciesCount int  `json:"dependencies_count"`
	IssuesCompsCount  int  `json:"issues_comps_count"`
	Modules           []struct {
		ModuleId       int    `json:"module_id"`
		Language       string `json:"language"`
		PackageManager string `json:"package_manager"`
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
	TaskId string `json:"task_id"`
	Status string `json:"status"`
}
