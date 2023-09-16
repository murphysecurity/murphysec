package api

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/version"
	"time"
)

type CreateSubTaskRequest struct {
	AccessType  model.AccessType `json:"access_type"` // 接入方式
	ScanMode    model.ScanMode   `json:"scan_mode"`
	Addr        *string          `json:"addr,omitempty"`
	Author      *string          `json:"author,omitempty"` // 作者/提交者Email
	Branch      *string          `json:"branch,omitempty"`
	Commit      *string          `json:"commit,omitempty"`
	Message     *string          `json:"message,omitempty"`   // 提交信息
	PushTime    *time.Time       `json:"push_time,omitempty"` // 提交时间
	SubtaskName string           `json:"subtask_name"`
	Dir         string           `json:"dir"` // 路径
	CliVersion  string           `json:"cli_version"`
	IsBuild     bool             `json:"is_build"`
	IsDeep      bool             `json:"is_deep"`
	ProjectName string           `json:"project_name"`
}

type CreateSubTaskResponse struct {
	ProjectsName string `json:"projects_name"` // 项目名称
	TaskID       string `json:"task_id"`       // 任务ID
	SubtaskID    string `json:"subtask_id"`    // 子任务ID
	TaskName     string `json:"task_name"`     // 任务名称
	AlertMessage string `json:"alert_message"`
}

func CreateSubTask(client *Client, request *CreateSubTaskRequest) (*CreateSubTaskResponse, error) {
	checkNotNull(client)
	checkNotNull(request)
	request.CliVersion = version.Version()
	var resp CreateSubTaskResponse
	if e := client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/create_subtask"), request), &resp); e != nil {
		return nil, e
	}
	return &resp, nil
}
