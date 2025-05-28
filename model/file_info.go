package model

// ProcessInfo 进程相关信息
type ProcessInfo struct {
	ProcessName string `json:"process_name,omitempty"`
	PID         string `json:"pid,omitempty"`
}

// FileInfo 文件相关信息
type FileInfo struct {
	FilePath    string      `json:"file_path,omitempty"`
	FileHash    string      `json:"file_hash,omitempty"`
	ProcessInfo ProcessInfo `json:"process_info,omitempty"`
}
