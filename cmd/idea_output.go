package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/utils"
	"murphysec-cli-simple/utils/must"
)

func reportIdeaErr(e error, message string) {
	code := IdeaUnknownErr
	if errors.Is(e, api.ErrTokenInvalid) {
		code = IdeaTokenInvalid
	} else if errors.Is(e, api.ErrServerRequest) {
		code = IdeaServerRequestFailed
	} else if errors.Is(e, api.ErrTimeout) {
		code = IdeaApiTimeout
	} else if errors.Is(e, api.BaseCommonApiError) {
		code = IdeaServerRequestFailed
	}
	if message == "" {
		message = e.Error()
	}
	if message == "" {
		message = code.Error()
	}
	fmt.Println(string(must.A(json.Marshal(struct {
		ErrCode IdeaErrCode `json:"err_code"`
		ErrMsg  string      `json:"err_msg"`
	}{ErrCode: code, ErrMsg: message}))))
}

type PluginOutput struct {
	ErrCode          IdeaErrCode  `json:"err_code"`
	IssuesCount      int          `json:"issues_count,omitempty"`
	Comps            []PluginComp `json:"comps,omitempty"`
	IssuesLevelCount struct {
		Critical int `json:"critical,omitempty"`
		High     int `json:"high,omitempty"`
		Medium   int `json:"medium,omitempty"`
		Low      int `json:"low,omitempty"`
	} `json:"issues_level_count,omitempty"`
	TaskId            string               `json:"task_id,omitempty"`
	TotalContributors int                  `json:"total_contributors"`
	ProjectId         string               `json:"project_id"`
	InspectErrors     []model.InspectError `json:"inspect_errors,omitempty"`
	DependenciesCount int                  `json:"dependencies_count"`
}
type PluginComp struct {
	CompName           string               `json:"comp_name"`
	ShowLevel          int                  `json:"show_level"`
	MinFixedVersion    string               `json:"min_fixed_version"`
	MinFixed           PluginCompFixList    `json:"min_fixed"`
	Vulns              []model.VoVulnInfo   `json:"vulns"`
	Version            string               `json:"version"`
	License            *PluginCompLicense   `json:"license,omitempty"`
	Solutions          []PluginCompSolution `json:"solutions,omitempty"`
	IsDirectDependency bool                 `json:"is_direct_dependency"`
	Language           string               `json:"language"`
	FixType            string               `json:"fix_type"`
}

type PluginCompLicense struct {
	Level model.LicenseLevel `json:"level"`
	Spdx  string             `json:"spdx"`
}

type PluginCompFix struct {
	OldVersion string `json:"old_version"`
	NewVersion string `json:"new_version"`
	CompName   string `json:"comp_name"`
}

type PluginCompFixList []PluginCompFix

func (this PluginCompFixList) MarshalJSON() ([]byte, error) {
	m := map[PluginCompFix]struct{}{}
	for _, it := range this {
		m[it] = struct{}{}
	}
	rs := make([]PluginCompFix, 0)
	for it := range m {
		rs = append(rs, it)
	}
	return must.A(json.Marshal(rs)), nil
}

type PluginCompSolution struct {
	Compatibility *int   `json:"compatibility,omitempty"`
	Description   string `json:"description"`
	Type          string `json:"type,omitempty"`
}

func generatePluginOutput(c context.Context) *PluginOutput {
	ctx := model.UseScanTask(c)
	i := ctx.ScanResult
	type id struct {
		name    string
		version string
	}
	p := &PluginOutput{
		ErrCode:     0,
		IssuesCount: i.IssuesCompsCount,
		Comps:       []PluginComp{},
		TaskId:      i.TaskId,
		//InspectErrors:     ctx.InspectorError,
		TotalContributors: ctx.TotalContributors,
		ProjectId:         ctx.ProjectId,
		DependenciesCount: ctx.ScanResult.DependenciesCount,
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
				MinFixed:        PluginCompFixList{},
				Vulns:           comp.Vuls,
				Version:         comp.CompVersion,
				License:         nil,
				Solutions:       []PluginCompSolution{},
				Language:        mod.Language,
				FixType:         comp.FixType,
			}
			for _, it := range comp.MinFixedInfo {
				p.MinFixed = append(p.MinFixed, PluginCompFix{
					OldVersion: it.OldVersion,
					NewVersion: it.NewVersion,
					CompName:   it.Name,
				})
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
				case model.SuggestLevelRecommend:
					p.ShowLevel = utils.MinInt(p.ShowLevel, 2)
				case model.SuggestLevelStrongRecommend:
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
					case model.VulnLevelCritical:
						critical[vul.VulnNo] = struct{}{}
					case model.VulnLevelHigh:
						high[vul.VulnNo] = struct{}{}
					case model.VulnLevelMedium:
						medium[vul.VulnNo] = struct{}{}
					case model.VulnLevelLow:
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
