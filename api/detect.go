package api

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
	Dependencies       interface{} `json:"dependencies"`
	RuntimeInfo        interface{} `json:"runtime_info"`
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
	Poc         bool                  `json:"poc"`
	Cvss        float64               `json:"cvss"`
	Description string                `json:"description"`
	Solution    interface{}           `json:"solution"`
	Source      string                `json:"source"`
	Effect      []ScanResultEffect    `json:"effect"`
	Suggest     string                `json:"suggest"`
	CompName    string                `json:"comp_name"`
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
