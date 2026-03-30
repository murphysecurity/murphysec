package model

// SkillSummary 表示 Skills 检测摘要。
type SkillSummary struct {
	Enabled             bool `json:"enabled"`
	Total               int  `json:"total"`
	RiskCount           int  `json:"risk_count"`
	MaliciousCount      int  `json:"malicious_count"`
	SuspiciousCount     int  `json:"suspicious_count"`
	SafeCount           int  `json:"safe_count"`
	PendingCount        int  `json:"pending_count"`
	DetectionErrorCount int  `json:"detection_error_count"`
	FailedCount         int  `json:"failed_count"`
	DetectorAvailable   bool `json:"detector_available"`
	LLMAvailable        bool `json:"llm_available"`
}

// SkillItem 表示单个 Skill 的检测结果。
type SkillItem struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	DirPath       string                `json:"dir_path"`
	DirName       string                `json:"dir_name"`
	Status        string                `json:"status"`
	Summary       string                `json:"summary"`
	DetectMethods []string              `json:"detect_methods"`
	Source        SkillSourceInfo       `json:"source"`
	Intelligence  SkillIntelligenceInfo `json:"intelligence"`
	Analysis      SkillAnalysisInfo     `json:"analysis"`
	ErrorMessage  string                `json:"error_message"`
}

// SkillSourceInfo 表示来源信息。
type SkillSourceInfo struct {
	SourceMatched bool   `json:"source_matched"`
	Platform      string `json:"platform"`
	PlatformURL   string `json:"platform_url"`
	Author        string `json:"author"`
	Version       string `json:"version"`
	Stars         int64  `json:"stars"`
	OriginalName  string `json:"original_name"`
}

// SkillIntelligenceInfo 表示情报匹配结果。
type SkillIntelligenceInfo struct {
	IsMalicious bool   `json:"is_malicious"`
	MpsID       string `json:"mps_id"`
	Summary     string `json:"summary"`
	CollectedAt string `json:"collected_at"`
}

// SkillAnalysisInfo 表示 LLM 分析结果。
type SkillAnalysisInfo struct {
	Analyzed      bool                     `json:"analyzed"`
	ModelProvider string                   `json:"model_provider"`
	ModelName     string                   `json:"model_name"`
	Summary       string                   `json:"summary"`
	Dimensions    []SkillAnalysisDimension `json:"dimensions"`
}

// SkillAnalysisDimension 表示单个分析维度结果。
type SkillAnalysisDimension struct {
	Name               string `json:"name"`
	DisplayName        string `json:"display_name"`
	Status             string `json:"status"`
	Detail             string `json:"detail"`
	Evidence           string `json:"evidence"`
	Location           string `json:"location"`
	UserExplanation    string `json:"user_explanation"`
	CommonRiskPatterns string `json:"common_risk_patterns"`
}
