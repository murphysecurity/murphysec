package inspector

import (
	"encoding/json"
	"fmt"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
	"time"
)

type IdeaErrCode int

const (
	IdeaSucceed          IdeaErrCode = 0
	IdeaEngineScanFailed IdeaErrCode = 1
	IdeaAPIFailed        IdeaErrCode = 2
	IdeaNoEngineMatch    IdeaErrCode = 3
	IdeaTokenInvalid     IdeaErrCode = 4
	IdeaUnknownErr       IdeaErrCode = -1
)

func reportIdeaStatus(code IdeaErrCode, msg string) {
	fmt.Println(string(must.Byte(json.Marshal(PluginOutput{ErrCode: code, ErrMsg: msg}))))
}

func mapForIdea(i *api.TaskScanResponse, ctx *ScanContext) PluginOutput {
	type id struct {
		name    string
		version string
	}
	p := PluginOutput{
		IssuesCount:            i.IssuesCompsCount,
		DependenciesCount:      i.DependenciesCount,
		IssuesCompsCount:       i.IssuesCompsCount,
		Comps:                  []PluginComp{},
		DetectorStartTimestamp: i.DetectStartTimestamp,
		DetectStatus:           i.DetectStatus,
		TaskId:                 i.TaskId,
	}
	// merge module comps
	rs := map[id]PluginComp{}
	for _, mod := range i.Modules {
		for _, comp := range mod.Comps {
			cid := id{comp.CompName, comp.CompVersion}
			p := PluginComp{
				CompName:        comp.CompName,
				ShowLevel:       3,
				MinFixedVersion: comp.MinFixedVersion,
				Vulns:           comp.Vuls,
				Version:         comp.CompVersion,
				License:         nil,
				Solutions:       []PluginCompSolution{},
				Language:        mod.Language,
			}
			if comp.License != nil {
				p.License = &PluginCompLicense{
					Level: comp.License.Level,
					Spdx:  comp.License.Spdx,
				}
			}
			// Work-around to keep result consistency.
			if rs[cid].IsDirectDependency {
				p.IsDirectDependency = true
			} else {
				p.IsDirectDependency = comp.IsDirectDependency
			}
			for _, it := range comp.Solutions {
				p.Solutions = append(p.Solutions, PluginCompSolution{
					Compatibility: it.Compatibility,
					Description:   it.Description,
					Type:          it.Type,
				})
			}
			for _, it := range comp.Vuls {
				switch it.SuggestLevel {
				case api.SuggestLevelRecommend:
					p.ShowLevel = utils.MinInt(p.ShowLevel, 2)
				case api.SuggestLevelStrongRecommend:
					p.ShowLevel = utils.MinInt(p.ShowLevel, 1)
				}
			}
			rs[cid] = p
		}
	}
	for _, it := range rs {
		p.Comps = append(p.Comps, it)
	}
	// calc vulns
	{
		critical := map[string]struct{}{}
		high := map[string]struct{}{}
		medium := map[string]struct{}{}
		low := map[string]struct{}{}
		for _, it := range i.Modules {
			for _, comp := range it.Comps {
				for _, vul := range comp.Vuls {
					switch vul.Level {
					case api.VulnLevelCritical:
						critical[vul.VulnNo] = struct{}{}
					case api.VulnLevelHigh:
						high[vul.VulnNo] = struct{}{}
					case api.VulnLevelMedium:
						medium[vul.VulnNo] = struct{}{}
					case api.VulnLevelLow:
						low[vul.VulnNo] = struct{}{}
					}
				}
			}
		}
		p.IssuesLevelCount.Low = len(low)
		p.IssuesLevelCount.Medium = len(medium)
		p.IssuesLevelCount.High = len(high)
		p.IssuesLevelCount.Critical = len(critical)
	}
	return p
}

type PluginOutput struct {
	ErrCode                IdeaErrCode  `json:"err_code"`
	ErrMsg                 string       `json:"err_msg,omitempty"`
	IssuesCount            int          `json:"issues_count,omitempty"`
	DependenciesCount      int          `json:"dependencies_count,omitempty"`
	IssuesCompsCount       int          `json:"issues_comps_count,omitempty"`
	Comps                  []PluginComp `json:"comps,omitempty"`
	DetectorStartTimestamp time.Time    `json:"detector_start_timestamp,omitempty"`
	DetectStatus           string       `json:"detect_status,omitempty"`
	IssuesLevelCount       struct {
		Critical int `json:"critical,omitempty"`
		High     int `json:"high,omitempty"`
		Medium   int `json:"medium,omitempty"`
		Low      int `json:"low,omitempty"`
	} `json:"issues_level_count,omitempty"`
	TaskId string `json:"task_id,omitempty"`
}

func (p PluginOutput) MarshalJSON() ([]byte, error) {
	if p.ErrCode != 0 {
		return must.Byte(json.Marshal(map[string]interface{}{"err_code": p.ErrCode, "err_msg": p.ErrMsg})), nil
	}
	type t PluginOutput
	return must.Byte(json.Marshal(t(p))), nil
}

type PluginComp struct {
	CompName           string               `json:"comp_name"`
	ShowLevel          int                  `json:"show_level"`
	MinFixedVersion    string               `json:"min_fixed_version"`
	Vulns              []api.VoVulnInfo     `json:"vulns"`
	Version            string               `json:"version"`
	License            *PluginCompLicense   `json:"license,omitempty"`
	Solutions          []PluginCompSolution `json:"solutions,omitempty"`
	IsDirectDependency bool                 `json:"is_direct_dependency"`
	Language           string               `json:"language"`
}

type PluginCompLicense struct {
	Level api.LicenseLevel `json:"level"`
	Spdx  string           `json:"spdx"`
}

type PluginCompSolution struct {
	Compatibility *int   `json:"compatibility,omitempty"`
	Description   string `json:"description"`
	Type          string `json:"type,omitempty"`
}
