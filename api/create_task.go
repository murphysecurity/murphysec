package api

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"murphysec-cli-simple/logger"
	"murphysec-cli-simple/utils/must"
	"net/http"
)

type CreateTaskRequest struct {
	CliVersion    string          `json:"cli_version"`
	TaskType      InspectTaskType `json:"task_type"`
	UserAgent     string          `json:"user_agent"`
	CmdLine       string          `json:"cmd_line"`
	ApiToken      string          `json:"api_token"`
	GitInfo       *VoGitInfo      `json:"git_info,omitempty"`
	ProjectName   string          `json:"project_name"`
	TargetAbsPath string          `json:"target_abs_path"`
}

type CreateTaskGitInfo struct {
	GitRemoteUrl string `json:"git_remote_url"`
	GitBranch    string `json:"git_branch"`
	Commit       string `json:"commit"`
}

type CreateTaskResponse struct {
	TaskInfo string `json:"task_info"`
}

func CreateTask(req *CreateTaskRequest) (*string, error) {
	body := must.Byte(json.Marshal(req))
	resp, e := http.Post(serverAddress()+"/message/v2/access/client/create_project", "application/json", bytes.NewReader(body))
	if e != nil {
		logger.Err.Println("Request failed", e.Error())
		return nil, errors.Wrap(ErrSendRequest, e.Error())
	}
	data, e := readHttpBody(resp)
	if e != nil {
		return nil, e
	}
	if resp.StatusCode == http.StatusOK {
		type O struct {
			Data CreateTaskResponse `json:"data"`
		}
		var o O
		if e := json.Unmarshal(data, &o); e != nil {
			return nil, e
		}
		if o.Data.TaskInfo == "" {
			return nil, errors.New("empty task info")
		}
		return &o.Data.TaskInfo, nil
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrTokenInvalid
	}
	return nil, readCommonErr(data, resp.StatusCode)
}
