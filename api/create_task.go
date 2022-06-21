package api

import (
	"github.com/murphysecurity/murphysec/model"
)

type CreateTaskRequest struct {
	CliVersion      string              `json:"cli_version"`
	TaskKind        model.TaskKind      `json:"task_kind"`
	TaskType        model.TaskType      `json:"task_type"`
	UserAgent       string              `json:"user_agent"`
	CmdLine         string              `json:"cmd_line"`
	ApiToken        string              `json:"api_token"`
	GitInfo         *VoGitInfo          `json:"git_info,omitempty"`
	ProjectName     string              `json:"project_name"`
	TargetAbsPath   string              `json:"target_abs_path"`
	ProjectType     model.ProjectType   `json:"project_type"`
	ContributorList []model.Contributor `json:"contributor_list,omitempty"`
	ProjectId       string              `json:"project_id,omitempty"`
}
type Contributor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateTaskGitInfo struct {
	GitRemoteUrl string `json:"git_remote_url"`
	GitBranch    string `json:"git_branch"`
	Commit       string `json:"commit"`
}

type CreateTaskResponse struct {
	TaskInfo          string `json:"task_info"`
	TotalContributors int    `json:"total_contributors"`
	ProjectId         string `json:"project_id"`
	Username          string `json:"username"`
}

func CreateTask(req *CreateTaskRequest) (*CreateTaskResponse, error) {
	httpReq := C.PostJson("/message/v2/access/client/create_project", req)
	type O struct {
		Data CreateTaskResponse `json:"data"`
	}
	var resp O
	if e := C.DoJson(httpReq, &resp); e != nil {
		return nil, e
	}
	return &resp.Data, nil
}
