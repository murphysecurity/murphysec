package api

import (
	"github.com/google/uuid"
	"murphysec-cli-simple/model"
	"murphysec-cli-simple/utils/must"
	"time"
)

type ModuleType string

const (
	ModuleTypeVersion  ModuleType = "version"
	ModuleTypeFileHash ModuleType = "file_hash"
)

type VoGitInfo struct {
	Commit        string    `json:"commit"`
	GitRef        string    `json:"git_ref"`
	GitRemoteUrl  string    `json:"git_remote_url"`
	CommitMessage string    `json:"commit_message"`
	CommitEmail   string    `json:"commit_email"`
	CommitTime    time.Time `json:"commit_time"`
}

type VoFileHash struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}
type VoModule struct {
	Dependencies   []model.Dependency       `json:"dependencies,omitempty"`
	FileHashList   []VoFileHash             `json:"file_hash_list,omitempty"`
	Language       model.Language           `json:"language,omitempty"`
	Name           string                   `json:"name,omitempty"`
	PackageFile    string                   `json:"package_file,omitempty"`
	PackageManager model.PackageManagerType `json:"package_manager,omitempty"`
	RelativePath   string                   `json:"relative_path,omitempty"`
	RuntimeInfo    interface{}              `json:"runtime_info,omitempty"`
	Version        string                   `json:"version,omitempty"`
	ModuleUUID     uuid.UUID                `json:"module_uuid,omitempty"`
	ModuleType     ModuleType               `json:"module_type"`
}

type SendDetectRequest struct {
	TaskInfo string     `json:"task_info"`
	ApiToken string     `json:"api_token"`
	Modules  []VoModule `json:"modules"`
}

func SendDetect(input *SendDetectRequest) error {
	must.True(input != nil)
	input.ApiToken = C.Token
	req := C.PostJson("/message/v2/access/detect/user_cli", input)
	return C.DoJson(req, nil)
}
