package inspector

import (
	"encoding/json"
	"fmt"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
	"time"
)

func mapForIdea(i *api.VoDetectResponse) PluginOutput {
	type id struct {
		name    string
		version string
	}
	p := PluginOutput{
		IssuesCount:            i.IssuesCompsCount,
		DependenciesCount:      i.DependenciesCount,
		IssuesCompsCount:       i.IssuesCompsCount,
		Language:               "java",
		PackageManager:         "maven",
		Comps:                  []PluginComp{},
		DetectorStartTimestamp: i.DetectStartTimestamp,
		DetectStatus:           i.DetectStatus,
		TaskId:                 i.TaskId,
	}
	// merge module comps
	rs := map[id]PluginComp{}
	for _, mod := range i.Modules {
		for _, comp := range mod.Comps {
			if _, ok := rs[id{comp.CompName, comp.CompVersion}]; ok {
				continue
			}
			p := PluginComp{
				CompName:        comp.CompName,
				ShowLevel:       3,
				MinFixedVersion: comp.MinFixedVersion,
				Vulns:           comp.Vuls,
			}
			for _, it := range comp.Vuls {
				switch it.SuggestLevel {
				case api.SuggestLevelRecommend:
					p.ShowLevel = utils.MinInt(p.ShowLevel, 2)
				case api.SuggestLevelStrongRecommend:
					p.ShowLevel = utils.MinInt(p.ShowLevel, 1)
				}
			}
			rs[id{comp.CompName, comp.CompVersion}] = p
		}
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
func ideaFail(code int, message string) {
	fmt.Println(string(must.Byte(json.Marshal(PluginOutput{ErrCode: code, ErrMsg: message}))))
}

type PluginOutput struct {
	ErrCode                int          `json:"err_code"`
	ErrMsg                 string       `json:"err_msg,omitempty"`
	IssuesCount            int          `json:"issues_count,omitempty"`
	DependenciesCount      int          `json:"dependencies_count,omitempty"`
	IssuesCompsCount       int          `json:"issues_comps_count,omitempty"`
	Language               string       `json:"language,omitempty"`
	PackageManager         string       `json:"package_manager,omitempty"`
	Comps                  []PluginComp `json:"comps,omitempty"`
	DetectorStartTimestamp time.Time    `json:"detector_start_timestamp,omitempty"`
	DetectStatus           string       `json:"detect_status,omitempty"`
	IssuesLevelCount       struct {
		Critical int `json:"critical"`
		High     int `json:"high"`
		Medium   int `json:"medium"`
		Low      int `json:"low"`
	} `json:"issues_level_count,omitempty"`
	TaskId string `json:"task_id,omitempty"`
}

type PluginComp struct {
	CompName        string           `json:"comp_name"`
	ShowLevel       int              `json:"show_level"`
	MinFixedVersion string           `json:"min_fixed_version"`
	Vulns           []api.VoVulnInfo `json:"vulns"`
}
