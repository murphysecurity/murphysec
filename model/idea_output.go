package model

import (
	"encoding/json"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
)

func GenerateIdeaErrorOutput(e error) string {
	iec := GetIdeaErrCode(e)
	return string(must.A(json.Marshal(struct {
		ErrCode IdeaErrCode `json:"err_code"`
		ErrMsg  string      `json:"err_msg"`
	}{ErrCode: iec, ErrMsg: e.Error()})))
}

type PluginOutput struct {
	SubtaskName      string       `json:"subtask_name"`
	ErrCode          IdeaErrCode  `json:"err_code"`
	IssuesCount      int          `json:"issues_count,omitempty"`
	Comps            []PluginComp `json:"comps,omitempty"`
	IssuesLevelCount struct {
		Critical int `json:"critical,omitempty"`
		High     int `json:"high,omitempty"`
		Medium   int `json:"medium,omitempty"`
		Low      int `json:"low,omitempty"`
	} `json:"issues_level_count,omitempty"`
	TaskId            string         `json:"task_id"`
	SubtaskId         string         `json:"subtask_id"`
	InspectErrors     []InspectError `json:"inspect_errors,omitempty"`
	DependenciesCount int            `json:"dependencies_count"`
	SurpassScore      int            `json:"surpass_score"`
	ProjectScore      int            `json:"project_score"`
}

type PluginComp struct {
	CompName           string                 `json:"comp_name"`
	ShowLevel          int                    `json:"show_level"`
	MinFixedVersion    string                 `json:"min_fixed_version"`
	Vulns              []PluginVulnDetailInfo `json:"vulns"`
	CompVersion        string                 `json:"version"`
	Licenses           []LicenseItem          `json:"licenses"`
	Solutions          []Solution             `json:"solutions,omitempty"`
	IsDirectDependency bool                   `json:"is_direct_dependency"`
	CompSecScore       int                    `json:"comp_sec_score"`
	FixPlanList        FixPlanList            `json:"fix_plan_list"`
	DependentPath      []string               `json:"dependent_path"`
	PackageManager     string                 `json:"package_manager"`
}

type PluginCompSolution struct {
	Compatibility *int   `json:"compatibility,omitempty"`
	Description   string `json:"description"`
	Type          string `json:"type,omitempty"`
}

func GetIDEAOutput(task *ScanTask) PluginOutput {

	// workaround: 从模块列表里拎包管理器出来
	pmMap := make(map[Component]string)
	for _, module := range task.Modules {
		for _, component := range module.ComponentList() {
			pmMap[component] = module.PackageManager
		}
	}
	var r = task.Result
	var pluginOutput = PluginOutput{
		ErrCode:     100,
		SubtaskName: task.SubtaskName,
		IssuesCount: r.LeakNum,
		Comps:       make([]PluginComp, 0),
		IssuesLevelCount: struct {
			Critical int `json:"critical,omitempty"`
			High     int `json:"high,omitempty"`
			Medium   int `json:"medium,omitempty"`
			Low      int `json:"low,omitempty"`
		}{
			Critical: r.CriticalNum,
			High:     r.HighNum,
			Medium:   r.MediumNum,
			Low:      r.LowNum,
		},
		TaskId:            r.TaskId,
		SubtaskId:         r.SubtaskId,
		InspectErrors:     nil,
		DependenciesCount: r.RelyNum,
		SurpassScore:      r.SurpassScore,
		ProjectScore:      r.ProjectScore,
	}

	var vulnListMapper = func(effects []ScanResultCompEffect) (rs []PluginVulnDetailInfo) {
		for _, effect := range effects {
			info, ok := r.VulnInfoMap[effect.MpsId]
			if !ok {
				continue // skip item if detailed information not found
			}
			var d = PluginVulnDetailInfo{
				MpsId:           info.MpsID,
				CveId:           info.CveID,
				Description:     info.Description,
				Level:           info.Level,
				Influence:       info.Influence,
				Poc:             info.Poc,
				PublishTime:     int(info.PublishedDate.Unix()),
				AffectedVersion: effect.EffectVersion,
				MinFixedVersion: effect.MinFixedVersion,
				References:      utils.NoNilSlice(info.ReferenceURLList),
				Solutions:       utils.NoNilSlice(effect.Solutions),
				SuggestLevel:    info.FixSuggestionLevel,
				Title:           info.Title,
			}
			rs = append(rs, d)
		}
		return
	}

	for _, comp := range r.CompInfoList {
		var pc = PluginComp{
			CompName:           comp.CompName,
			CompVersion:        comp.CompVersion,
			ShowLevel:          comp.SuggestLevel,
			MinFixedVersion:    comp.MinFixedVersion,
			Vulns:              utils.NoNilSlice(vulnListMapper(comp.VulnList)),
			Licenses:           utils.NoNilSlice(comp.LicenseList),
			Solutions:          utils.NoNilSlice(comp.Solutions),
			IsDirectDependency: comp.IdDirectDependency,
			CompSecScore:       comp.CompSecScore,
			FixPlanList:        comp.FixPlanList,
			DependentPath:      utils.NoNilSlice(comp.DependentPath),
			PackageManager:     pmMap[comp.Component],
		}
		if len(pc.Vulns) == 0 {
			continue
		}
		// workaround: IDE侧要求我一定加进去，后续他不要求了，就删掉
		if pc.FixPlanList.Plan1 != nil {
			pc.FixPlanList.Plan1.CompName = comp.CompName
			pc.FixPlanList.Plan1.OldVersion = comp.CompVersion
		}
		if pc.FixPlanList.Plan2 != nil {
			pc.FixPlanList.Plan2.CompName = comp.CompName
			pc.FixPlanList.Plan2.OldVersion = comp.CompVersion
		}
		if pc.FixPlanList.Plan3 != nil {
			pc.FixPlanList.Plan3.CompName = comp.CompName
			pc.FixPlanList.Plan3.OldVersion = comp.CompVersion
		}

		pluginOutput.Comps = append(pluginOutput.Comps, pc)
	}
	return pluginOutput
}

type PluginVulnDetailInfo struct {
	MpsId           string         `json:"mps_id"`
	CveId           string         `json:"cve_id"`
	Description     string         `json:"description"`
	Level           string         `json:"level"`
	Influence       int            `json:"influence"`
	Poc             bool           `json:"poc"`
	PublishTime     int            `json:"publish_time"`
	AffectedVersion string         `json:"affected_version"`
	MinFixedVersion string         `json:"min_fixed_version"`
	References      []ReferenceURL `json:"references"`
	Solutions       []Solution     `json:"solutions"`
	SuggestLevel    string         `json:"suggest_level"`
	Title           string         `json:"title"`
}
