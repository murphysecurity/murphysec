package model

import "time"

// ScanResultResponse 服务端输出，字段名比较迷惑，需要注意。
//
// 标注为 <unknown> 的暂时不知道是什么，估计用不上
type ScanResultResponse struct {
	Complete     bool                 `json:"complete"`
	SubtaskId    string               `json:"subtask_id"`
	ProjectsId   string               `json:"projects_id"`
	TeamId       string               `json:"team_id"`
	TaskId       string               `json:"task_id"`
	UserId       string               `json:"user_id"`
	FinishTime   time.Time            `json:"finish_time"`
	Addr         string               `json:"addr"`          // <unknown>
	FromSource   string               `json:"from_source"`   // <unknown>
	OptionalNum  int                  `json:"optional_num"`  // 可选修复数量
	RecommendNum int                  `json:"recommend_num"` // 建议修复数量
	StringNum    int                  `json:"string_num"`    // 强烈建议修复数量
	RelyNum      int                  `json:"rely_num"`      // 本次扫描依赖组件总数
	LeakNum      int                  `json:"leak_num"`      // 本次任务包含漏洞总数
	HighNum      int                  `json:"high_num"`      // 高危漏洞数量
	MediumNum    int                  `json:"medium_num"`    // 中危漏洞数量
	LowNum       int                  `json:"low_num"`       // 低危漏洞数量
	CriticalNum  int                  `json:"critical_num"`  // 严重漏洞数量
	LicenseNum   int                  `json:"license_num"`   // 许可证数量
	Languages    string               `json:"languages"`     // 语言，逗号隔开的
	SurpassScore int                  `json:"surpass_score"`
	ProjectScore int                  `json:"project_score"`
	CompInfoList []ScanResultCompInfo `json:"comp_info_list"`
	VulnInfoMap  map[string]VulnerabilityDetailInfo
}

type ScanResultCompInfo struct {
	Component
	IdDirectDependency bool                   `json:"id_direct_dependency"`
	CompSecScore       int                    `json:"comp_sec_score"`
	MinFixedVersion    string                 `json:"min_fixed_version"`
	CriticalNum        int                    `json:"critical_num"`
	HighNum            int                    `json:"high_num"`
	MediumNum          int                    `json:"medium_num"`
	LowNum             int                    `json:"low_num"`
	VulnList           []ScanResultCompEffect `json:"vuln_list,omitempty"`
	LicenseList        []LicenseItem          `json:"license_list,omitempty"`
	DependentPath      []string               `json:"dependent_path,omitempty"`
	Solutions          []Solution             `json:"solutions,omitempty"`
	FixPlanList        FixPlanList            `json:"fix_plan_list"`
	SuggestLevel       int                    `json:"suggest_level"` // 对应到IDEA的show_level，具体计算规则不明
}

type ScanResultCompEffect struct {
	EffectVersion   string     `json:"effect_version"`
	MinFixedVersion string     `json:"min_fixed_version"`
	MpsId           string     `json:"mps_id"`
	Solutions       []Solution `json:"solutions,omitempty"`
}

type Solution struct {
	Description     string `json:"description"`
	Type            string `json:"type"`
	CompatibleScore *int   `json:"compatible_score,omitempty"`
}

type LicenseItem struct {
	Spdx  string       `json:"spdx"`
	Level LicenseLevel `json:"level"`
}

type LicenseLevel string

const (
	LicenseLevelHigh   LicenseLevel = "High"
	LicenseLevelMedium LicenseLevel = "Medium"
	LicenseLevelLow    LicenseLevel = "Low"
)

type FixPlanItem struct {
	CompatibilityScore int    `json:"compatibility_score"`
	SecurityScore      int    `json:"security_score"`
	TargetVersion      string `json:"target_version"`
	CompName           string `json:"comp_name,omitempty"` // 这两个字段现在不应该放在这，但是IDE那边一定要我塞进去
	OldVersion         string `json:"old_version,omitempty"`
}

type FixPlanList struct {
	Plan1 *FixPlanItem `json:"plan1,omitempty"`
	Plan2 *FixPlanItem `json:"plan2,omitempty"`
	Plan3 *FixPlanItem `json:"plan3,omitempty"`
}

// VulnerabilityDetailInfo 漏洞详情
type VulnerabilityDetailInfo struct {
	AttackVector       string         `json:"attack_vector"`        // 攻击向量
	CnvdID             string         `json:"cnvd_id"`              // 漏洞CNVD ID
	CveID              string         `json:"cve_id"`               // 漏洞CVE ID
	CvssScore          float64        `json:"cvss_score"`           //
	CvssVector         string         `json:"cvss_vector"`          // CVSS 向量
	Description        string         `json:"description"`          // 漏洞详情信息
	Exp                bool           `json:"exp"`                  // 是否有EXP
	Exploitability     string         `json:"exploitability"`       //
	FixSuggestionLevel string         `json:"fix_suggestion_level"` //
	Influence          int            `json:"influence"`            // 漏洞影响指数
	Languages          []string       `json:"languages"`            // 漏洞语言
	Level              string         `json:"level"`                //
	MpsID              string         `json:"mps_id"`               // 漏洞MPS ID
	Patch              string         `json:"patch"`                // Patch信息
	Poc                bool           `json:"poc"`                  // 存在POC与否？
	PublishedDate      time.Time      `json:"published_date"`       // 漏洞发布时间
	ReferenceURLList   []ReferenceURL `json:"reference_url_list"`   //
	ScopeInfluence     string         `json:"scope_influence"`      //
	Title              string         `json:"title"`                // 漏洞标题
	TroubleShooting    []string       `json:"trouble_shooting"`     // 排查方式列表
	VulnType           string         `json:"vuln_type"`            // 漏洞类型
}

type ReferenceURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
